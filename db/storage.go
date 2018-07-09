package db

import "github.com/tsuru/cst/scan"

// Storage represents a persistent data store.
type Storage interface {
	Close()
	Save(scan.Scan) error
}
