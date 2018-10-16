package passes

import (
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
	"time"
)

// DAO is the data access object for interacting with TrackingPass structs
// stored in MongoDB
type DAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	// COLLECTION is the MongoDB collection in which TrackingPass structs are stored
	COLLECTION = "passes"
)

// Connect connects the PassesDAO to a MongoDB server
func (d *DAO) Connect() {
	session, err := mgo.Dial(d.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(d.Database)
}

// FindAll retrieves all TrackingPass object from MongoDB and returns them in
// order by start date/time
func (d *DAO) FindAll() ([]TrackingPass, error) {
	var passes []TrackingPass
	err := db.C(COLLECTION).Find(bson.M{}).Sort("start_time").All(&passes)
	return passes, err
}

// Insert adds a TrackingPass to MongoDB
func (d *DAO) Insert(pass TrackingPass) error {
	// TODO: check to see if pass already exists
	err := db.C(COLLECTION).Insert(&pass)
	return err
}

// FindByID finds and retrieves a TrackingPass struct based upon its ID
func (d *DAO) FindByID(id string) (TrackingPass, error) {
	var pass TrackingPass
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&pass)
	return pass, err
}

// FindByQuery finds TrackingPass structs associated with a given
// MongoDB query (see https://docs.mongodb.com/manual/tutorial/query-documents/
// for additional information)
func (d *DAO) FindByQuery(query bson.M) ([]TrackingPass, error) {
	var passes []TrackingPass
	err := db.C(COLLECTION).Find(query).Sort("start_time").All(&passes)
	return passes, err
}

// GetNextPass retrieves the next TrackingPass (the closest in the future)
func (d *DAO) GetNextPass() (TrackingPass, error) {
	t := time.Now()
	var passes []TrackingPass
	err := db.C(COLLECTION).Find(bson.M{"start_time": bson.M{"$gte": t}}).Sort("start_time").All(&passes)
	return passes[0], err
}

// Delete removes a TrackingPass from MongoDB
func (d *DAO) Delete(pass TrackingPass) error {
	err := db.C(COLLECTION).Remove(&pass)
	return err
}

// Update edits the info for a TrackingPass in MongoDB
func (d *DAO) Update(pass TrackingPass) error {
	err := db.C(COLLECTION).UpdateId(pass.ID, &pass)
	return err
}
