package internals

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var logsFolder string = "logs/"

type Database struct {
	store map[string]string
	mutex *sync.RWMutex
}

type FastStoreCache interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	BuildAtTimeOfStart() error
}

func NewDatabase() FastStoreCache {
	return &Database{
		store: map[string]string{},
		mutex: &sync.RWMutex{},
	}
}

// Set store key,value in map
func (d *Database) Set(key, value string) error {
	// add this to log file
	data := fmt.Sprintf("%s|%s|%s", "set", key, value)

	now := time.Now()
	dateInt, err := strconv.Atoi(now.Format("20060102"))
	if err != nil {
		log.Fatal(err)
	}
	currentTime := fmt.Sprintf("%d-%02d-%02d", dateInt, now.Hour(), now.Minute())
	currentFileName := fmt.Sprintf("%s-ops.wal", currentTime)

	err = AppendLineToFile(currentFileName, data)
	if err != nil {
		return err
	}
	d.mutex.Lock()
	d.store[key] = value
	d.mutex.Unlock()
	return nil
}

// Get retrieve value of particular key
func (d *Database) Get(key string) (string, error) {
	data := fmt.Sprintf("%s|%s", "get", key)

	now := time.Now()
	dateInt, err := strconv.Atoi(now.Format("20060102"))
	if err != nil {
		log.Fatal(err)
	}
	currentTime := fmt.Sprintf("%d-%02d-%02d", dateInt, now.Hour(), now.Minute())
	currentFileName := fmt.Sprintf("%s-ops.wal", currentTime)

	err = AppendLineToFile(currentFileName, data)
	if err != nil {
		return "", err
	}
	d.mutex.Lock()
	val, ok := d.store[key]
	d.mutex.Unlock()
	if !ok {
		return "", nil
	}
	return val, nil
}

// Delete remove particular key
func (d *Database) Delete(key string) error {
	data := fmt.Sprintf("%s|%s", "del", key)
	now := time.Now()
	dateInt, err := strconv.Atoi(now.Format("20060102"))
	if err != nil {
		log.Fatal(err)
	}
	currentTime := fmt.Sprintf("%d-%02d-%02d", dateInt, now.Hour(), now.Minute())
	currentFileName := fmt.Sprintf("%s-ops.wal", currentTime)

	err = AppendLineToFile(currentFileName, data)
	if err != nil {
		return err
	}
	d.mutex.Lock()
	delete(d.store, key)
	d.mutex.Unlock()
	return nil
}

// BuildAtTimeOfStart function will build the db.This is basically for data persistence
func (d *Database) BuildAtTimeOfStart() error {
	file, err := OpenFile("logs/global.wal")
	if err != nil {
		return err
	}
	log.Println("FileName is ", file.Name())
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println("inside this")
		line := scanner.Text()
		cmds := strings.Split(line, "|")
		switch cmds[0] {
		case "set":
			d.store[cmds[1]] = cmds[2]

		case "del":
			delete(d.store, cmds[1])

		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Inside thisssssss ", err)
		return err
	}
	return nil
}
