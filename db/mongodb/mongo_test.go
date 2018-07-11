package mongodb

import (
	"testing"

	"github.com/globalsign/mgo"

	"github.com/stretchr/testify/assert"
	"github.com/tsuru/cst/scan"

	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
}

func TestMongoDB_Save(t *testing.T) {

	if !viper.IsSet("STORAGE_URL") {
		t.Skip("mongodb connection url are not assigned, skipping integration tests")
	}

	mongo, err := NewMongoDB(viper.GetString("STORAGE_URL"))

	if err != nil {
		assert.FailNow(t, "could not connect with mongodb")
	}

	defer func() {
		scanColl := mongo.getScanCollection()
		scanColl.Database.Session.Close()

		mongo.session.Close()
	}()

	t.Run(`When a scan document already exists on datastore, should only update that document`, func(t *testing.T) {

		scanColl := mongo.getScanCollection()

		defer func() {
			scanColl.DropCollection()
			scanColl.Database.Session.Close()
		}()

		scanColl.Insert(scan.Scan{
			ID:     "2b935a8f-4241-49f0-a1a2-e3c8ba347b95",
			Image:  "tsuru/cst:latest",
			Status: scan.StatusScheduled,
		})

		mongo.Save(scan.Scan{
			ID:     "2b935a8f-4241-49f0-a1a2-e3c8ba347b95",
			Image:  "tsuru/cst:latest",
			Status: scan.StatusFinished,
		})

		var got scan.Scan

		scanColl.FindId("2b935a8f-4241-49f0-a1a2-e3c8ba347b95").One(&got)

		expected := scan.Scan{
			ID:     "2b935a8f-4241-49f0-a1a2-e3c8ba347b95",
			Image:  "tsuru/cst:latest",
			Status: scan.StatusFinished,
		}

		assert.Equal(t, expected, got)
	})

	t.Run(`When a scan document not exists on datastore, should inserts it`, func(t *testing.T) {

		scanColl := mongo.getScanCollection()

		defer func() {
			scanColl.DropCollection()
			scanColl.Database.Session.Close()
		}()

		newScan := scan.Scan{
			ID:     "83633447-353f-4e87-aa95-2a44205eb89e",
			Image:  "tsuru/cst:latest",
			Status: scan.StatusScheduled,
		}

		var got scan.Scan

		err := scanColl.FindId(newScan.ID).One(&got)

		assert.Error(t, err)
		assert.Equal(t, mgo.ErrNotFound, err)

		mongo.Save(newScan)

		scanColl.FindId(newScan.ID).One(&got)

		assert.Equal(t, newScan, got)
	})
}

func TestMongoDB_Close(t *testing.T) {

	if !viper.IsSet("STORAGE_URL") {
		t.Skip("mongodb connection url are not assigned, skipping integration tests")
	}

	mongo, err := NewMongoDB(viper.GetString("STORAGE_URL"))

	if err != nil {
		assert.FailNow(t, "could not connect with mongodb")
	}

	t.Run(`Ensure any command issued after MongoDB.Close should panic the execution`, func(t *testing.T) {

		assert.Panics(t, func() {
			mongo.Close()

			mongo.session.Ping()
		})
	})
}

func TestMongoDB_HasScheduledScanByImage(t *testing.T) {

	if !viper.IsSet("STORAGE_URL") {
		t.Skip("mongodb connection url are not assigned, skipping integration tests")
	}

	mongo, err := NewMongoDB(viper.GetString("STORAGE_URL"))

	if err != nil {
		assert.FailNow(t, "could not connect with mongodb")
	}

	t.Run(`When exists a scan document with same image and status scheduled, should return true`, func(t *testing.T) {

		scanColl := mongo.getScanCollection()

		defer func() {
			scanColl.DropCollection()
			scanColl.Database.Session.Close()
		}()

		scanColl.Insert(scan.Scan{
			Image:  "tsuru/cst:latest",
			Status: scan.StatusScheduled,
		})

		assert.True(t, mongo.HasScheduledScanByImage("tsuru/cst:latest"))
	})

	t.Run(`When exist a scan document but with status non-scheduled, should return false`, func(t *testing.T) {

		scanColl := mongo.getScanCollection()

		defer func() {
			scanColl.DropCollection()
			scanColl.Database.Session.Close()
		}()

		scanColl.Insert(scan.Scan{
			ID:     "2b935a8f-4241-49f0-a1a2-e3c8ba347b95",
			Image:  "tsuru/cst:latest",
			Status: scan.StatusFinished,
		})

		assert.False(t, mongo.HasScheduledScanByImage("tsuru/cst:latest"))
	})
}

func TestMongoDB_AppendResultToScanByID(t *testing.T) {

	if !viper.IsSet("STORAGE_URL") {
		t.Skip("mongodb connection url are not assigned, skipping integration tests")
	}

	mongo, err := NewMongoDB(viper.GetString("STORAGE_URL"))

	if err != nil {
		assert.FailNow(t, "could not connect with mongodb")
	}

	t.Run(`When a scan has no results yet, should return one result after`, func(t *testing.T) {

		scanColl := mongo.getScanCollection()

		defer func() {
			scanColl.DropCollection()
			scanColl.Database.Session.Close()
		}()

		scanColl.Insert(scan.Scan{
			ID:     "2b935a8f-4241-49f0-a1a2-e3c8ba347b95",
			Image:  "tsuru/cst:latest",
			Result: []scan.Result{},
		})

		var scanOnStorage scan.Scan

		scanColl.FindId("2b935a8f-4241-49f0-a1a2-e3c8ba347b95").One(&scanOnStorage)

		assert.Equal(t, 0, len(scanOnStorage.Result))

		err := mongo.AppendResultToScanByID("2b935a8f-4241-49f0-a1a2-e3c8ba347b95", scan.Result{
			Scanner:         "scanner-example",
			Vulnerabilities: "all-vulns-described-here",
		})

		if assert.NoError(t, err) {
			scanColl.FindId("2b935a8f-4241-49f0-a1a2-e3c8ba347b95").One(&scanOnStorage)

			assert.Equal(t, 1, len(scanOnStorage.Result))
		}
	})
}

func TestMongoDB_UpdateScanStatusByID(t *testing.T) {

	if !viper.IsSet("STORAGE_URL") {
		t.Skip("mongodb connection url are not assigned, skipping integration tests")
	}

	mongo, err := NewMongoDB(viper.GetString("STORAGE_URL"))

	if err != nil {
		assert.FailNow(t, "could not connect with mongodb")
	}

	t.Run(`When updating scan to running status, should update scan status to running`, func(t *testing.T) {

		scanColl := mongo.getScanCollection()

		defer func() {
			scanColl.DropCollection()
			scanColl.Database.Session.Close()
		}()

		scanColl.Insert(scan.Scan{
			ID:     "2b935a8f-4241-49f0-a1a2-e3c8ba347b95",
			Status: scan.StatusScheduled,
		})

		var scanOnStorage scan.Scan

		scanColl.FindId("2b935a8f-4241-49f0-a1a2-e3c8ba347b95").One(&scanOnStorage)

		assert.Equal(t, scan.StatusScheduled, scanOnStorage.Status)

		err := mongo.UpdateScanStatusByID("2b935a8f-4241-49f0-a1a2-e3c8ba347b95", scan.StatusRunning)

		if assert.NoError(t, err) {
			scanColl.FindId("2b935a8f-4241-49f0-a1a2-e3c8ba347b95").One(&scanOnStorage)

			assert.Equal(t, scan.StatusRunning, scanOnStorage.Status)
		}
	})
}
