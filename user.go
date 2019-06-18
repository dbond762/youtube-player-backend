package youtube_player_backend

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64
	Login    string
	Password string
	Likes    []Video
}

func (u *User) CheckPassword(password []byte) (bool, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	if u.Password == string(hash) {
		return true, nil
	}
	return false, nil
}

func HashPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

type UserService interface {
	UserByID(id int) (*User, error)
	UserByLogin(login string) (*User, error)
	CreateUser(u *User) (int64, error)
}
