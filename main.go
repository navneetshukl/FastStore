package main

import (
	"context"
	"fast-store/internals"
	"fast-store/server"
	"flag"
	"fmt"
	"time"
)

func main() {
	port := flag.String("port", "6369", "faststore server port")
	userName := flag.String("username", "", "username for db server")
	password := flag.String("password", "", "password for db server")

	flag.Parse()
	cacheService := internals.NewDatabase()

	err := cacheService.BuildAtTimeOfStart()
	fmt.Println("Err from read is ", err)

	newDBServer := server.NewServer(*userName, *port, *password, cacheService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTimer(70 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Server Closing")
				return

			case <-ticker.C:
				internals.LogCompaction()
			}
		}
	}()

	newDBServer.StartServer()

}
