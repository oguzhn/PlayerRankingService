package business

import (
	"errors"
	"log"

	"github.com/oguzhn/PlayerRankingService/database"
	"github.com/oguzhn/PlayerRankingService/models"
)

type Business struct {
	handler database.IDatabase
}

func NewBusiness(handler database.IDatabase) *Business {
	return &Business{handler: handler}
}

func (b *Business) AddScore(score models.ScoreDTO) error {
	return nil
}

func (b *Business) CreateUser(user models.UserDTO) error {
	_, err := b.handler.GetUserById(user.ID) //check if its id already exists on db before creating new user
	if err == nil {
		return errors.New("Same id already created for another user")
	}
	root, err := b.handler.GetRoot()
	if err != nil {
		log.Println(err)
		return err
	}
	if root == nil {
		err = b.handler.SetRoot(user.ID)
		if err != nil {
			log.Println(err)
			return err
		}
		userdb := migrateModelToDbModel(user)
		userdb.RightCount = 1
		return b.handler.CreateUser(&userdb)
	}

	p := root
	for p != nil {
		if user.Score >= p.Score {
			p.RightCount++
			err = b.handler.UpdateUser(p)
			if err != nil {
				log.Println(err)
				return err
			}
			if p.RightID != "" {
				p, err = b.handler.GetUserById(p.RightID)
				if err != nil {
					log.Println(err)
					return err
				}
			} else {
				p.RightID = user.ID
				err = b.handler.UpdateUser(p)
				if err != nil {
					log.Println(err)
					return err
				}
				userdb := migrateModelToDbModel(user)
				err = b.handler.CreateUser(&userdb)
				if err != nil {
					log.Println(err)
					return err
				}
				break
			}
		} else {
			if p.LeftID != "" {
				p, err = b.handler.GetUserById(p.LeftID)
				if err != nil {
					log.Println(err)
					return err
				}
			} else {
				p.LeftID = user.ID
				err = b.handler.UpdateUser(p)
				if err != nil {
					log.Println(err)
					return err
				}
				userdb := migrateModelToDbModel(user)
				err = b.handler.CreateUser(&userdb)
				if err != nil {
					log.Println(err)
					return err
				}
				break
			}
		}
	}

	return nil
}

func (b *Business) GetLeaderBoardByCountryCode(countryCode string) (models.UserDTOList, error) {
	leaderboard, err := b.GetLeaderBoard()
	if err != nil {
		return leaderboard, err
	}
	return leaderboard.FilterByCountryCode(countryCode), nil
}

func (b *Business) GetLeaderBoard() (models.UserDTOList, error) {
	var resultList models.UserDTOList
	root, err := b.handler.GetRoot()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resultList, err = b.bstTraverse(root)
	for i := range resultList {
		resultList[i].Rank = i + 1
	}
	return resultList, err
}

func (b *Business) bstTraverse(root *database.UserDAO) (models.UserDTOList, error) {
	resultList := models.UserDTOList{}

	if root.RightID != "" {
		rightUser, err := b.handler.GetUserById(root.RightID)
		if err != nil {
			return resultList, err
		}
		rightList, err := b.bstTraverse(rightUser)
		if err != nil {
			return resultList, err
		}
		resultList = append(resultList, rightList...)
	}
	resultList = append(resultList, migrateDbModelToModel(*root))
	if root.LeftID != "" {
		leftUser, err := b.handler.GetUserById(root.LeftID)
		if err != nil {
			return resultList, err
		}
		leftList, err := b.bstTraverse(leftUser)
		if err != nil {
			return resultList, err
		}
		resultList = append(resultList, leftList...)
	}
	return resultList, nil
}

func (b *Business) GetUserById(id string) (models.UserDTO, error) {
	user, err := b.handler.GetUserById(id)
	if err != nil {
		log.Println(err)
		return models.UserDTO{}, err
	}
	userDTO := migrateDbModelToModel(*user)

	rank := 0
	root, err := b.handler.GetRoot()
	if err != nil {
		log.Println(err)
		return models.UserDTO{}, err
	}

	p := root
	for p.ID != userDTO.ID {
		if userDTO.Score >= p.Score {
			p, err = b.handler.GetUserById(p.RightID)
			if err != nil {
				log.Println(err)
				return userDTO, err
			}
		} else {
			rank += p.RightCount
			p, err = b.handler.GetUserById(p.LeftID)
			if err != nil {
				log.Println(err)
				return userDTO, err
			}
		}
	}
	rank += p.RightCount
	userDTO.Rank = rank
	return userDTO, nil
}
