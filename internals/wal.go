package internals

import (
	"bufio"
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

// read file line by line
func ReadFromFile() error {
	file, err := os.Open("ops.wal")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text() // one line at a time
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
