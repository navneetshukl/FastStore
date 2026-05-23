package internals

type Database struct {
	store map[string]string
}

func (d *Database) Set(key, value string) bool {
	d.store[key] = value
	return true
}

func (d *Database) Get(key string) string {
	if _, ok := d.store[key]; !ok {
		return ""
	}
	return d.store[key]
}

func (d *Database) Delete(key string) bool {
	delete(d.store, key)
	return true
}
