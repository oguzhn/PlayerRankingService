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

func (cl *MongoClient) UpdateUser(u *UserDAO) error {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)

	return collection.UpdateId(u.Id, u)
}

func (cl *MongoClient) CreateUser(u *UserDAO) error {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	return collection.Insert(u)
}

func (cl *MongoClient) GetUserById(id string) (*UserDAO, error) {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	var user *UserDAO
	err := collection.FindId(id).One(user)
	return user, err
}

func (cl *MongoClient) GetRoot() (*UserDAO, error) {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	var root *UserDAO
	err := collection.FindId(rootId).One(root)
	if err != nil {
		return nil, err
	}

	err = collection.FindId(root.LeftId).One(root)

	return root, err
}

func (cl *MongoClient) RemoveUserById(id string) error {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	return collection.RemoveId(id)
}
