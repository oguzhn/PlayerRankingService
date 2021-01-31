type IDatabase interface {
	UpdateUser(u *UserDAO) error
	CreateUser(u *UserDAO) error
	GetUserById(id string) (*UserDAO, error)
	GetRoot() (*UserDAO, error)
	RemoveUserById(id string) error
}