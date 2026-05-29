package internals

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// log will be stored like ops|key|value (eg -> SET|name|navneet)
// create a new log file every minute and store the log in this file.
func WriteToFile(data string) error {

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return err
	}

	now := time.Now()
	dateInt, err := strconv.Atoi(now.Format("20060102"))
	if err != nil {
		log.Fatal(err)
	}
	currentTime := fmt.Sprintf("%d-%02d-%02d", dateInt, now.Hour(), now.Minute())

	// every minute new log file is created with format date|hour|minute
	logFileName := fmt.Sprintf("logs/%s-ops.wal", currentTime)
	file, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
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

func LogCompaction(){

	// read all the file present in logs aprt from global.wal
}