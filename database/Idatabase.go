package database

type IDatabase interface {
	UpdateUser(u *UserDAO) error
	CreateUser(u *UserDAO) error
	GetUserById(id string) (*UserDAO, error)
	GetRoot() (*UserDAO, error)
	SetRoot(id string) error
	RemoveUserById(id string) error
}
