package business

import (
	"github.com/oguzhn/PlayerRankingService/database"
	"github.com/oguzhn/PlayerRankingService/models"
)

func migrateModelToDbModel(user models.UserDTO) database.UserDAO {
	return database.UserDAO{
		ID:          user.ID,
		Score:       user.Score,
		Name:        user.Name,
		CountryCode: user.CountryCode,
		LeftID:      "",
		RightID:     "",
		RightCount:  1,
	}
}

func migrateDbModelToModel(user database.UserDAO) models.UserDTO {
	return models.UserDTO{
		ID:          user.ID,
		Score:       user.Score,
		Name:        user.Name,
		CountryCode: user.CountryCode,
	}
}
