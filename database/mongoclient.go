package database

import "github.com/globalsign/mgo"

type MongoClient struct {
	ms      *mgo.Session
	dbName  string
	colName string
}

// NewDatastore create connection to mongo db given the connection string
func NewDatastore(con string, db string, collection string) (*MongoClient, error) {
	ms, err := mgo.Dial(con)
	if err != nil {
		return nil, err
	}
	return &MongoClient{ms, db, collection}, nil
}
