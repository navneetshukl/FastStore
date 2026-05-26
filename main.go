package main

import (
	"fast-store/internals"
	"fmt"
)

func main() {
	db := internals.NewDatabase()

	err:=db.BuildAtTimeOfStart()
	fmt.Println("Err from read is ",err)

	// for i := 0; i < 5; i++ {
	// 	db.Set(fmt.Sprintf("idx%d",i), fmt.Sprintf("%d",i))
	// 	//db.Get("shukla")
	// }

	//  db.Delete("name")

	val,_:=db.Get("idx1")
	fmt.Println("Val is ",val)

	//internals.ReadFromFile()

}
