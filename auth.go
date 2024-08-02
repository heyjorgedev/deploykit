package deploykit

import "context"

type Registration struct {
	Name  string
	Email string
}

type RegistrationService struct {
	userService UserService
}

func NewRegistrationService(userService UserService) *RegistrationService {
	return &RegistrationService{userService: userService}
}

func (s *RegistrationService) Register(ctx context.Context, registration *Registration) error {
	user := &User{
		Name:  registration.Name,
		Email: registration.Email,
	}

	return s.userService.CreateUser(ctx, user)
}
