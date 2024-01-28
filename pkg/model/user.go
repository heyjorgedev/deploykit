package model

type User struct {
	ID           int64  `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"passwordHash" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u User) ValidatePassword(password string) bool {
	// TODO: Change
	return password == u.PasswordHash
}
