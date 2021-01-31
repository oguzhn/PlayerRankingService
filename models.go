package main

type UserDTO struct {
	Id          string  `json:"user_id"`
	Score       float32 `json:"points"`
	Name        string  `json:"display_name"`
	Rank        int     `json:"rank"`
	CountryCode string  `json:"country"`
}

type ScoreDTO struct {
	Score     float32 `json:"score_worth"`
	Id        string  `json:"user_id"`
	Timestamp int     `json:"timestamp"`
}
