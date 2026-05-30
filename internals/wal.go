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

// CreateFolderIfNotExist create folder if it not exists
func CreateFolderIfNotExist(folderName string) error {
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return err
	}
	return nil
}

// OpenFile create file if it does not exist
func OpenFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// AppendLineToFile add single line to file
func AppendLineToFile(fileName string, data string) error {

	file, err := OpenFile(fileName)
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

// AppendBulkData write bulk data to file
func AppendBulkData(fileName string, bulkData string) error {
	file, err := OpenFile(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(bulkData)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

// log will be stored like ops|key|value (eg -> SET|name|navneet)

// read file line by line
func ReadFromFile(fileName string) (*os.File, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
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
			var builder strings.Builder
			for k, v := range logsData {
				data := fmt.Sprintf("set|%s|%s\n", k, v)
				builder.WriteString(data)

			}
			err = AppendBulkData("logs/global.wal", builder.String())
			if err != nil {
				log.Println("error in writing bulk to file")
				return
			}

			// move this log file to processed log file
		}
		file.Close()
		//return
	}
}
