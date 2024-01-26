package mock

import "github.com/heyjorgedev/deploykit"

type UserService struct {
	FindAllFunc        func() (*[]deploykit.User, error)
	CreateFunc         func(user *deploykit.User) error
	DeleteFunc         func(user *deploykit.User) error
	FindByIdFunc       func(id int) (*deploykit.User, error)
	FindByUsernameFunc func(username string) (*deploykit.User, error)
}

func (u *UserService) FindAll() (*[]deploykit.User, error) {
	return u.FindAllFunc()
}

func (u *UserService) Create(user *deploykit.User) error {
	return u.CreateFunc(user)
}

func (u *UserService) Delete(user *deploykit.User) error {
	return u.DeleteFunc(user)
}

func (u *UserService) FindById(id int) (*deploykit.User, error) {
	return u.FindByIdFunc(id)
}

func (u *UserService) FindByUsername(username string) (*deploykit.User, error) {
	return u.FindByUsernameFunc(username)
}
