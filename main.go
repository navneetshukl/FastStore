package main

import (
	"fast-store/internals"
	"fast-store/server"
	"flag"
	"fmt"
)

func main() {
	port := flag.String("port", "6369", "faststore server port")
	userName := flag.String("username", "", "username for db server")
	password := flag.String("password", "", "password for db server")

	flag.Parse()
	db := internals.NewDatabase()

	err := db.BuildAtTimeOfStart()
	fmt.Println("Err from read is ", err)

	// for i := 0; i < 100; i++ {
	// 	db.Set(fmt.Sprintf("idx%d", i), fmt.Sprintf("%d", i))
	// 	//db.Get("shukla")
	// 	log.Println("SET : ", i)
	// 	time.Sleep(1 * time.Second)
	// }

	//  db.Delete("name")

	// val, _ := db.Get("idx1")
	// fmt.Println("Val is ", val)

	//internals.ReadFromFile()

	newDBServer := server.NewServer(*userName, *port, *password)

	newDBServer.StartServer()

	//internals.LogCompaction()

}
