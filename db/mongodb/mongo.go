package mongodb

import (
	"github.com/globalsign/mgo"
	"github.com/tsuru/cst/scan"
)

// MongoDB implements a Storage interface.
type MongoDB struct {
	session *mgo.Session
}

// Save inserts or updates (if scan.ID already exists on current collection)
// a scan document on MongoDB service.
func (mongo *MongoDB) Save(scan scan.Scan) error {

	collection := mongo.getScanCollection()
	defer collection.Database.Session.Close()

	_, err := collection.UpsertId(scan.ID, scan)

	return err
}

// Close permanently terminates the session with MongoDB service.
func (mongo *MongoDB) Close() {
	mongo.session.Close()
}

func (mongo *MongoDB) getScanCollection() *mgo.Collection {

	session := mongo.session.Copy()

	return session.DB("").C("scans")
}

// NewMongoDB creates a new instance of MongoDB and estabilishes a new session
// with MongoDB service. Returns an error if MongoDB service is unavailable.
func NewMongoDB(rawURL string) (*MongoDB, error) {

	session, err := mgo.Dial(rawURL)

	if err != nil {
		return nil, err
	}

	return &MongoDB{
		session: session,
	}, nil
}
