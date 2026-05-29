package main

import (
	"fast-store/internals"
	"fmt"
)

func main() {
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

	// server.StartServer()

	internals.LogCompaction()

}
