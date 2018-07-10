package db

import "github.com/tsuru/cst/scan"

// Storage represents a persistent data store.
type Storage interface {
	Close()
	Save(scan.Scan) error
}

var storageInstance Storage

// SetStorage sets the storage parameter to be used globally by any that calls
// GetStorage function.
func SetStorage(storage Storage) {
	storageInstance = storage
}

// GetStorage returns the current storage instance. Make sure you set a storage
// instance (calling the SetStorage function) before.
func GetStorage() Storage {
	return storageInstance
}
