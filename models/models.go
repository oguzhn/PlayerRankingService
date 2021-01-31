package models

import (
	"errors"

	"github.com/beevik/guid"
)

type UserDTO struct {
	Id          string  `json:"user_id"`
	Score       float32 `json:"points"`
	Name        string  `json:"display_name"`
	Rank        int     `json:"rank"`
	CountryCode string  `json:"country"`
}

type UserDTOList []UserDTO

type ScoreDTO struct {
	Score     float32 `json:"score_worth"`
	Id        string  `json:"user_id"`
	Timestamp int     `json:"timestamp"`
}

func (list UserDTOList) FilterByCountryCode(countryCode string) (ret UserDTOList) {
	for _, el := range list {
		if el.CountryCode == countryCode {
			ret = append(ret, el)
		}
	}
	return
}

func (s ScoreDTO) IsValid() error {

	if s.Score < 0 {
		return errors.New("Invalid score")
	}
	if s.Timestamp < 0 {
		return errors.New("Invalid Timestamp")
	}
	_, err := guid.ParseString(s.Id)

	return err
}

func (u UserDTO) IsValid() error {

	_, err := guid.ParseString(u.Id)

	return err
}
