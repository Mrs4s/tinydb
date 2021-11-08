package tinydb

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

// JSONStorage Create a new JSONStorage instance.
func JSONStorage(file string) (*StorageJSON, error) {
	var dir string
	i1 := strings.Index(file, `\`)
	i2 := strings.Index(file, `/`)
	if i1 != -1 || i2 != -1 {
		if i1 > i2 {
			dir = file[:i1]
		} else {
			dir = file[:i2]
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, err
			}
		}
	}
	fi, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	return &StorageJSON{handle: fi}, err
}

// Read Read data from JSON file
func (sto *StorageJSON) Read() (StorageData, error) {
	var data = StorageData{}
	sto.handle.Seek(0, 0)
	dec := json.NewDecoder(sto.handle)
	err := dec.Decode(&data)
	if err != nil && err != io.EOF {
		return data, err
	}
	return data, nil
}

// Write Write data to JSON file
func (sto *StorageJSON) Write(data StorageData) error {
	sto.mutex.Lock()
	defer sto.mutex.Unlock()
	sto.handle.Truncate(0)
	sto.handle.Seek(0, 0)
	enc := json.NewEncoder(sto.handle)
	enc.SetIndent("", "    ")
	if data == nil {
		return errors.New("Nothing needs to be written")
	}
	return enc.Encode(data)
}

// Close Close the JSONStorage instance.
func (sto *StorageJSON) Close() error {
	return sto.handle.Close()
}

// MemoryStorage Create a new MemoryStorage instance.
func MemoryStorage() (*StorageMemory, error) {
	return &StorageMemory{memory: StorageData{}}, nil
}

// Read Read data from memory
func (sto *StorageMemory) Read() (StorageData, error) {
	return sto.memory, nil
}

// Write Write data to memory
func (sto *StorageMemory) Write(data StorageData) error {
	sto.memory = data
	return nil
}

// MemoryStorage Close the MemoryStorage instance.
func (sto *StorageMemory) Close() error {
	return nil
}
