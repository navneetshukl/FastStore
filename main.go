package main

import "fast-store/internals"

func main(){
	f:=internals.Wal{
		FileName: "ops.wal",
	}

	f.WriteToFile("navneet")
	f.ReadFromFile()
}
