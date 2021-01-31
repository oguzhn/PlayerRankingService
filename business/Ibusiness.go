package business

import (
	"github.com/oguzhn/PlayerRankingService/models"
)

type IBusiness interface {
	AddScore(models.ScoreDTO) error
	CreateUser(models.UserDTO) error
	GetLeaderBoard() (models.UserDTOList, error)
	GetLeaderBoardByCountryCode(string) (models.UserDTOList, error)
	GetUserById(string) (models.UserDTO, error)
}
