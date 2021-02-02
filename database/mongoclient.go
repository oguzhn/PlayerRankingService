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

	return collection.UpdateId(u.ID, u)
}

func (cl *MongoClient) CreateUser(u *UserDAO) error {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	return collection.Insert(u)
}

func (cl *MongoClient) GetUserById(id string) (*UserDAO, error) {
	if id == "" {
		return nil, nil
	}
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	var user UserDAO
	err := collection.FindId(id).One(&user)
	return &user, err
}

func (cl *MongoClient) GetRoot() (*UserDAO, error) {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	var root UserDAO
	err := collection.FindId(rootId).One(&root)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil //root does not exist so it is not an error
		}
		return nil, err
	}

	err = collection.FindId(root.LeftID).One(&root) //leftId gives us the root user. I picked leftid randomly. It could have been right as well

	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil //root did not set so it is not an error
		}
		return nil, err
	}
	return &root, err
}

func (cl *MongoClient) SetRoot(id string) error {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)

	var root UserDAO
	err := collection.FindId(rootId).One(&root)
	if err != nil {
		if err != mgo.ErrNotFound {
			return err
		}
	}
	_, err = collection.UpsertId(rootId, &UserDAO{ID: rootId, LeftID: id})
	return err
}

func (cl *MongoClient) RemoveUserById(id string) error {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	return collection.RemoveId(id)
}
