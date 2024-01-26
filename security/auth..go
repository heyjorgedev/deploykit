package security

import (
	"github.com/heyjorgedev/deploykit"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userService deploykit.UserService
}

func NewAuthService(userService deploykit.UserService) *AuthService {
	return &AuthService{userService: userService}
}

func (a *AuthService) AttemptCredentials(username, password string) (*deploykit.User, error) {
	user, err := a.userService.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *AuthService) CreateAccount(user *deploykit.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return a.userService.Create(user)
}
