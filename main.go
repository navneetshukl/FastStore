package main

import "fast-store/internals"

func main() {
	db := internals.NewDatabase()

	for i := 0; i < 5; i++ {
		db.Set("name", "navneet")
		db.Get("shukla")
	}

}
