package mongodb

import (
	"github.com/globalsign/mgo"
	"github.com/tsuru/cst/scan"
	"gopkg.in/mgo.v2/bson"
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

// HasScheduledScanByImage checks if exists scan documents given image (by
// parameter) and status "scheduled" on MongoDB service. Returns true when
// exists one or more documents otherwise returns false.
func (mongo *MongoDB) HasScheduledScanByImage(image string) bool {

	collection := mongo.getScanCollection()
	defer collection.Database.Session.Close()

	scanFilter := scan.Scan{
		Image:  image,
		Status: scan.StatusScheduled,
	}

	documentsCount, _ := collection.Find(scanFilter).Count()

	return documentsCount > 0
}

// Close permanently terminates the session with MongoDB service.
func (mongo *MongoDB) Close() {
	mongo.session.Close()
}

// AppendResultToScanByID append the result on scan on MongoDB service.
func (mongo *MongoDB) AppendResultToScanByID(id string, result scan.Result) error {

	collection := mongo.getScanCollection()
	defer collection.Database.Session.Close()

	return collection.UpdateId(id, bson.M{"$push": bson.M{"result": result}})
}

// UpdateScanStatusByID updates the scan status on MongoDB service.
func (mongo *MongoDB) UpdateScanStatusByID(id string, status scan.Status) error {

	collection := mongo.getScanCollection()
	defer collection.Database.Session.Close()

	return collection.UpdateId(id, bson.M{"$set": bson.M{"status": status}})
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
