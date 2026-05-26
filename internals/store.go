package internals

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Database struct {
	store map[string]string
}

func NewDatabase() *Database {
	return &Database{
		store: map[string]string{},
	}
}

// Set store key,value in map
func (d *Database) Set(key, value string) error {
	// add this to log file
	data := fmt.Sprintf("%s|%s|%s", "set", key, value)
	err := WriteToFile(data)
	if err != nil {
		return err
	}
	d.store[key] = value
	return nil
}

// Get retrieve value of particular key
func (d *Database) Get(key string) (string, error) {
	data := fmt.Sprintf("%s|%s", "get", key)
	err := WriteToFile(data)
	if err != nil {
		return "", err
	}
	if _, ok := d.store[key]; !ok {
		return "", nil
	}
	return d.store[key], nil
}

// Delete remove particular key
func (d *Database) Delete(key string) error {
	data := fmt.Sprintf("%s|%s", "del", key)
	err := WriteToFile(data)
	if err != nil {
		return err
	}
	delete(d.store, key)
	return nil
}

// BuildAtTimeOfStart function will build the db.This is basically for data persistence
func (d *Database) BuildAtTimeOfStart()error {
	file, err := os.Open("ops.wal")
	if err != nil {
		return err
	}
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
		return err
	}
	return nil
}
