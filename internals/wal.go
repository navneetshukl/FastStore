package internals

import (
	"fmt"
	"os"
)

// log will be stored like ops|key|value (eg -> SET|name|navneet)
func WriteToFile(data string) error {
	file, err := os.OpenFile("ops.wal", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data + "\n")
	if err != nil {
		return err
	}
	return nil
}

func ReadFromFile() error {
	data, err := os.ReadFile("ops.wal")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
