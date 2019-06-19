package youtube_player_backend

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64
	Login    string
	Password string
}

func (u *User) CheckPassword(password []byte) bool {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error on check password: %s", err)
		return false
	}

	if u.Password == string(hash) {
		return true
	}
	return false
}

func HashPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error on gen hash: %s", err)
		return "", err
	}
	return string(hash), nil
}

type UserService interface {
	UserByID(id int) (*User, error)
	UserByLogin(login string) (*User, error)
	CreateUser(u User) (*User, error)
}

type UserSession interface {
	Login(login string) (*User, string, error)
	Logout(token string)
	Authenticate(token string) (*User, error)
}
