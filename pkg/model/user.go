package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID           int    `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"passwordHash" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	return nil
}
