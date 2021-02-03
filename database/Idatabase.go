package database

/*
IDatabase is the declaration of how the database behaves. It exposes this API to the upper layers.
*/
type IDatabase interface {
	UpdateUser(u *UserDAO) error
	CreateUser(u *UserDAO) error
	GetUserByID(id string) (*UserDAO, error)
	GetRoot() (*UserDAO, error)
	SetRoot(id string) error
	RemoveUserByID(id string) error
}
