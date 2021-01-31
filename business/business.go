package business

import "github.com/oguzhn/PlayerRankingService/models"

type Business struct {
}

func NewBusiness() *Business {
	return &Business{}
}

func (b *Business) AddScore(score models.ScoreDTO) error {
	return nil
}

func (b *Business) CreateUser(user models.UserDTO) error {
	return nil
}

func (b *Business) GetLeaderBoardByCountryCode(countryCode string) (models.UserDTOList, error) {
	return nil, nil
}

func (b *Business) GetLeaderBoard() (models.UserDTOList, error) {
	return nil, nil
}

func (b *Business) GetUserById(id string) (models.UserDTO, error) {

	return models.UserDTO{}, nil
}
