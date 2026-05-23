package internals

import (
	"fmt"
	"log"
	"os"
)

type Wal struct{
	FileName string
}

// log will be stored like ops|key|value (eg -> SET|name|navneet)
func (w *Wal)WriteToFile(data string) {
	file, err := os.OpenFile(w.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(data + "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func(w *Wal)ReadFromFile()error{
	data,err:=os.ReadFile(w.FileName)
	if err!=nil{
		return err
	}
	fmt.Println(string(data))
	return nil
}
