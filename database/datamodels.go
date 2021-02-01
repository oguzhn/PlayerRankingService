package database

const rootId string = "root"

type UserDAO struct {
	ID          string  `bson:"_id"`
	Score       float32 `bson:"score"`
	Name        string  `bson:"name"`
	CountryCode string  `bson:"country"`
	LeftID      string  `bson:"leftId"`
	RightID     string  `bson:"rightId"`
	RightCount  int     `bson:"rightCount"`
}
