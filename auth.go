package deploykit

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserService interface {
	FindAll() (*[]User, error)
	Create(user *User) error
	Delete(user *User) error
	FindById(id int) (*User, error)
	FindByUsername(username string) (*User, error)
}

type AuthService interface {
	AttemptCredentials(username, password string) (*User, error)
	CreateAccount(user *User) error
}
