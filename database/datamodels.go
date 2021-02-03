package database

//rootID keeps track of root node on mongodb
const rootID string = "root"

/*
UserDAO represents a user in the database. I implemented a Binary search tree that is hosted on a mongodb client.
That's why it keeps track of leftId and rightId of two different users. Left users score is always less than parent
users and right users is always greater than or equal to parent users. It also keeps track of rightcount which
represents the number of nodes that is greater than or equal to the given node. Since it is equal to itself, we include
itself in the rightcount. So a leaf node's rightcount is 1.
*/
type UserDAO struct {
	ID          string  `bson:"_id"`
	Score       float32 `bson:"score"`
	Name        string  `bson:"name"`
	CountryCode string  `bson:"country"`
	LeftID      string  `bson:"leftId"`
	RightID     string  `bson:"rightId"`
	RightCount  int     `bson:"rightCount"`
}
