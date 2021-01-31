package database

const rootId string = "root"

type UserDAO struct {
	Id          string  `bson:"_id"`
	Score       float32 `bson:"score"`
	Name        string  `bson:"name"`
	CountryCode string  `bson:"country"`
	LeftId      string  `bson:"leftId"`
	RightId     string  `bson:"rightId"`
	RightCount  int     `bson:"rightCount"`
}
