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
	user, err := b.handler.GetUserById(score.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	root, err := b.handler.GetRoot()
	if err != nil {
		log.Println("could not get root ", err)
		return err
	}
	if root.ID == user.ID {
		if root.RightID == "" {
			err = b.handler.SetRoot(root.LeftID)
			if err != nil {
				return err
			}
		} else if root.LeftID == "" {
			err = b.handler.SetRoot(root.RightID)
			if err != nil {
				return err
			}
		} else {
			err = b.handler.SetRoot(root.RightID)
			if err != nil {
				return err
			}
			userRight, err := b.handler.GetUserById(root.RightID)
			if err != nil {
				return err
			}
			left := userRight.LeftID
			userRight.LeftID = root.LeftID
			err = b.handler.UpdateUser(userRight)
			if err != nil {
				log.Println("could not update user ", err)
				return err
			}

			leftUser, err := b.handler.GetUserById(left)
			if err != nil {
				return err
			}

			rootsLeft, err := b.handler.GetUserById(root.LeftID)
			if err != nil {
				return err
			}

			count, err := b.countNumberOfNodes(leftUser)
			if err != nil {
				return err
			}

			for rootsLeft.RightID != "" {
				rootsLeft.RightCount += count
				err = b.handler.UpdateUser(rootsLeft)
				if err != nil {
					return err
				}
				rootsLeft, err = b.handler.GetUserById(rootsLeft.RightID)
				if err != nil {
					return err
				}
			}
			rootsLeft.RightCount += count
			rootsLeft.RightID = left
			err = b.handler.UpdateUser(rootsLeft)
			if err != nil {
				return err
			}
		}
	} else {
		parent, isLeftChild, err := b.findParentUser(user)
		if err != nil {
			return err
		}
		if user.RightID == "" {
			if isLeftChild {
				parent.LeftID = user.LeftID
			} else {
				parent.RightID = user.LeftID
				parent.RightCount--
				b.decrementUpperNodesRightCount(parent)
			}
			err = b.handler.UpdateUser(parent)
			if err != nil {
				return err
			}
		} else if user.LeftID == "" {
			if isLeftChild {
				parent.LeftID = user.RightID
			} else {
				parent.RightID = user.RightID
				parent.RightCount--
				b.decrementUpperNodesRightCount(parent)
			}
			err = b.handler.UpdateUser(parent)
			if err != nil {
				return err
			}
		} else {
			if !isLeftChild {
				parent.RightID = user.RightID
				parent.RightCount--
				b.decrementUpperNodesRightCount(parent)
			} else {
				parent.LeftID = user.RightID
			}
			err = b.handler.UpdateUser(parent)
			if err != nil {
				return err
			}
			userRight, err := b.handler.GetUserById(user.RightID)
			if err != nil {
				return err
			}
			left := userRight.LeftID
			userRight.LeftID = user.LeftID
			err = b.handler.UpdateUser(userRight)
			if err != nil {
				log.Println("could not update user ", err)
				return err
			}

			leftUser, err := b.handler.GetUserById(left)
			if err != nil {
				return err
			}

			usersLeft, err := b.handler.GetUserById(user.LeftID)
			if err != nil {
				return err
			}

			count, err := b.countNumberOfNodes(leftUser)
			if err != nil {
				return err
			}

			for usersLeft.RightID != "" {
				usersLeft.RightCount += count
				err = b.handler.UpdateUser(usersLeft)
				if err != nil {
					return err
				}
				usersLeft, err = b.handler.GetUserById(usersLeft.RightID)
				if err != nil {
					return err
				}
			}
			usersLeft.RightCount += count
			usersLeft.RightID = left
			err = b.handler.UpdateUser(usersLeft)
			if err != nil {
				return err
			}
		}
	}

	err = b.handler.RemoveUserById(user.ID)
	if err != nil {
		return err
	}
	return b.CreateUser(models.UserDTO{
		ID:          user.ID,
		Score:       score.Score,
		Name:        user.Name,
		CountryCode: user.CountryCode,
	})
}

func (b *Business) decrementUpperNodesRightCount(user *database.UserDAO) error {
	root, err := b.handler.GetRoot()
	if err != nil {
		return err
	}
	p := root
	for p.ID != user.ID {
		if user.Score >= p.Score {
			p.RightCount--
			b.handler.UpdateUser(p)
			p, err = b.handler.GetUserById(p.RightID)
			if err != nil {
				return err
			}
		} else {
			p, err = b.handler.GetUserById(p.LeftID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (b *Business) countNumberOfNodes(root *database.UserDAO) (int, error) {
	if root == nil || root.ID == "" {
		return 0, nil
	}

	var rightUser *database.UserDAO
	var err error
	if root.RightID != "" {
		rightUser, err = b.handler.GetUserById(root.RightID)
		if err != nil {
			return 0, err
		}
	}
	right, err := b.countNumberOfNodes(rightUser)
	if err != nil {
		return 0, err
	}

	var leftUser *database.UserDAO
	if root.LeftID != "" {
		leftUser, err = b.handler.GetUserById(root.LeftID)
		if err != nil {
			return 0, err
		}
	}
	left, err := b.countNumberOfNodes(leftUser)
	if err != nil {
		return 0, err
	}

	return right + left + 1, nil
}

func (b *Business) findParentUser(user *database.UserDAO) (*database.UserDAO, bool, error) { //bool is to determine if it is left or right of parent. return true if it is left
	root, err := b.handler.GetRoot()
	if err != nil {
		return nil, false, err
	}
	p := root
	for p != nil {
		if p.RightID == user.ID || p.LeftID == user.ID {
			return p, p.LeftID == user.ID, nil //return true if it is left
		}
		if user.Score >= p.Score {
			p, err = b.handler.GetUserById(p.RightID)
			if err != nil {
				return nil, false, err
			}
		} else {
			p, err = b.handler.GetUserById(p.LeftID)
			if err != nil {
				return nil, false, err
			}
		}
	}
	return nil, false, errors.New("could not find parent user")
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
