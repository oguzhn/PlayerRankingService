package models

import (
	"errors"

	"github.com/beevik/guid"
)

type UserDTO struct {
	ID          string  `json:"user_id"`
	Score       float32 `json:"points"`
	Name        string  `json:"display_name"`
	Rank        int     `json:"rank"`
	CountryCode string  `json:"country"`
}

type UserDTOList []UserDTO

type BulkUserDTO struct {
	Count int         `json:"count"`
	List  UserDTOList `json:"users"`
}

type ScoreDTO struct {
	Score     float32 `json:"score_worth"`
	ID        string  `json:"user_id"`
	Timestamp int     `json:"timestamp"`
}

func (list UserDTOList) FilterByCountryCode(countryCode string) (ret UserDTOList) {
	ret = UserDTOList{}
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
	_, err := guid.ParseString(s.ID)

	return err
}

func (u UserDTO) IsValid() error {

	_, err := guid.ParseString(u.ID)

	return err
}

func IsValidGuid(str string) error {
	_, err := guid.ParseString(str)

	return err
}
