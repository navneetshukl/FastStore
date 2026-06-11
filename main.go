package main

import (
	"context"
	"fast-store/internals"
	"fast-store/server"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	port := flag.String("port", "6369", "faststore server port")
	userName := flag.String("username", "", "username for db server")
	password := flag.String("password", "", "password for db server")

	flag.Parse()
	cacheService := internals.NewDatabase()

	err := cacheService.BuildAtTimeOfStart()
	log.Println("Err from read is ", err)

	newDBServer := server.NewServer(*userName, *port, *password, cacheService)

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Compaction worker stopped")
				return

			case <-ticker.C:
				log.Println("Starting log compaction")
				internals.LogCompaction()
				log.Println("Log compaction completed")
			}
		}
	}()

	go func() {
		newDBServer.StartServer()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %v", sig)

	cancel()

	log.Println("Waiting for background jobs to finish...")
	wg.Wait()

	log.Println("Graceful shutdown completed")
}
