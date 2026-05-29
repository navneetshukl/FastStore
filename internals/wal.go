package internals

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

//. TODO. : Refactor this file and break it in functions

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
func ReadFromFile(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	//defer file.Close()

	//scanner := bufio.NewScanner(file)

	// for scanner.Scan() {
	// 	line := scanner.Text() // one line at a time
	// 	fmt.Println(line)
	// }

	// if err := scanner.Err(); err != nil {
	// 	return err
	// }

	return file, nil
}

// LogCompaction function will compact the log file
func LogCompaction() {

	// read all the file present in logs aprt from global.wal
	files, err := os.ReadDir("logs")
	if err != nil {
		panic(err)
	}

	now := time.Now()
	dateInt, err := strconv.Atoi(now.Format("20060102"))
	if err != nil {
		log.Fatal(err)
	}
	currentTime := fmt.Sprintf("%d-%02d-%02d", dateInt, now.Hour(), now.Minute())

	// every minute new log file is created with format date|hour|minute
	currentFileName := fmt.Sprintf("%s-ops.wal", currentTime)

	fileNames := []string{}
	for _, file := range files {
		if file.Name() != "global.wal" || file.Name() != currentFileName {
			fileNames = append(fileNames, file.Name())
		}
	}

	//fmt.Println("Original file name is ", fileNames)
	sort.Slice(fileNames, func(i, j int) bool {
		return fileNames[i] < fileNames[j]
	})

	//fmt.Println("Sorted Order is ", fileNames)
	logsData := map[string]string{}

	for _, v := range fileNames {
		v = "logs/" + v
		file, err := ReadFromFile(v)
		if err != nil {
			log.Printf("Error in reading file %s is %v \n", v, err)
			continue
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			currentLine := scanner.Text()

			contents := strings.Split(currentLine, "|")
			if len(contents) <= 1 {
				continue
			}
			cmd := contents[0]
			key := contents[1]

			switch cmd {
			case "set":
				if len(contents) < 3 {
					break
				}
				logsData[key] = contents[2]

			case "del":
				delete(logsData, key)

			}
		}

		log.Println("logs Data is ", logsData)

		err = scanner.Err()
		fmt.Println("Error from read is ", v, " ", err)
		if err == nil {
			// write this map data to global.wal log file

			for k, v := range logsData {
				file, err := os.OpenFile("logs/global.wal", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
				if err != nil {
					fmt.Println("Error in opening the file ", err)
					return
				}
				defer file.Close()

				writer := bufio.NewWriter(file)

				data := fmt.Sprintf("set|%s|%s", k, v)

				_, err = writer.WriteString(data + "\n")
				if err != nil {
					fmt.Println("error in writer ", err)
					return
				}
				err = writer.Flush()
				if err != nil {
					fmt.Println("Error in flush ", err)
					return
				}

				err = file.Sync()
				if err != nil {
					fmt.Println("error in sync")
					return
				}
				logsData = map[string]string{}
			}

			// move this log file to processed log file
		}
		file.Close()
		//return
	}
}
