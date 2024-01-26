package security

import (
	"github.com/heyjorgedev/deploykit"
	"github.com/heyjorgedev/deploykit/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestAuthService_CreateAccount(t *testing.T) {
	createFuncCalled := false
	userService := &mock.UserService{
		CreateFunc: func(user *deploykit.User) error {
			createFuncCalled = true
			if user.Password == "password" {
				t.Errorf("expected password to be hashed")
			}

			if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password")) != nil {
				t.Errorf("expected hash comparison to succeed")
			}

			return nil
		},
	}
	a := AuthService{userService: userService}

	u := &deploykit.User{
		Name:     "User",
		Username: "example",
		Password: "password",
	}

	err := a.CreateAccount(u)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !createFuncCalled {
		t.Error("expected authService to delegate the creation of the user to the userService")
	}
}

func TestAuthService_AttemptCredentials(t *testing.T) {
	userService := &mock.UserService{
		FindByUsernameFunc: func(username string) (*deploykit.User, error) {
			u := &deploykit.User{
				Username: "example",
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			u.Password = string(hashedPassword)

			return u, nil
		},
	}

	a := AuthService{userService: userService}

	u, err := a.AttemptCredentials("example", "password")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if u.Username != "example" {
		t.Errorf("expected username to be example, got %s", u.Username)
	}
}
