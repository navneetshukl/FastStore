package internals

import "fmt"

type Database struct {
	store map[string]string
}

func NewDatabase()*Database{
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
