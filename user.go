package deploykit

import "context"

type User struct {
	ID    int
	Name  string
	Email string
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
}
