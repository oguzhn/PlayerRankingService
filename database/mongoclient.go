package database

import "github.com/globalsign/mgo"

/*MongoClient is a mongo client which has a session, database name and collection name*/
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

/*UpdateUser is simply updates a given user*/
func (cl *MongoClient) UpdateUser(u *UserDAO) error {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)

	return collection.UpdateId(u.ID, u)
}

/*CreateUser simply creates a new user*/
func (cl *MongoClient) CreateUser(u *UserDAO) error {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	return collection.Insert(u)
}

/*GetUserByID simply gets the user for given id*/
func (cl *MongoClient) GetUserByID(id string) (*UserDAO, error) {
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

/*
GetRoot get the root of Binary search tree. It first finds the id of the root
which is written on the mongo document with id specified in rootID constant variable.
Then it fetches the user with that id.
*/
func (cl *MongoClient) GetRoot() (*UserDAO, error) {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	var root UserDAO
	err := collection.FindId(rootID).One(&root)
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

/*
SetRoot sets the root for given id. It first checks if root exists then it sets it.
*/
func (cl *MongoClient) SetRoot(id string) error {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)

	var root UserDAO
	err := collection.FindId(rootID).One(&root)
	if err != nil {
		if err != mgo.ErrNotFound {
			return err
		}
	}
	_, err = collection.UpsertId(rootID, &UserDAO{ID: rootID, LeftID: id})
	return err
}

/*
RemoveUserByID removes the user with given id
*/
func (cl *MongoClient) RemoveUserByID(id string) error {
	session := cl.ms.Copy()
	defer session.Close()
	collection := session.DB(cl.dbName).C(cl.colName)
	return collection.RemoveId(id)
}
