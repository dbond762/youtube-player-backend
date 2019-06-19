package postgres

import (
	"database/sql"
	"log"

	player "github.com/dbond762/youtube-player-backend"
	_ "github.com/lib/pq"
)

type UserService struct {
	DB *sql.DB
}

func (s *UserService) UserByID(id int) (*player.User, error) {
	var u player.User
	row := s.DB.QueryRow(`SELECT "id", "login", "password" FROM "user" WHERE "id" = $1`, id)
	if err := row.Scan(&u.ID, &u.Login, &u.Password); err != nil {
		log.Printf("Postgres: Error on scan row: %s", err)
		return nil, err
	}

	return &u, nil
}

func (s *UserService) UserByLogin(login string) (*player.User, error) {
	var u player.User
	row := s.DB.QueryRow(`SELECT "id", "login", "password" FROM "user" WHERE "login" = $1`, login)
	if err := row.Scan(&u.ID, &u.Login, &u.Password); err != nil {
		log.Printf("Postgres: Error on scan row: %s", err)
		return nil, err
	}

	return &u, nil
}

func (s *UserService) CreateUser(u player.User) (*player.User, error) {
	var lastID int64
	hash, err := player.HashPassword([]byte(u.Password))
	if err != nil {
		log.Printf("Postgres: Error on hash password: %s", err)
		return nil, err
	}

	err = s.DB.QueryRow(`INSERT INTO "user" ("login", "password") VALUES ($1, $2)  RETURNING "id"`, u.Login, hash).Scan(&lastID)
	if err != nil {
		log.Printf("Postgres: Query error: %s", err)
		return nil, err
	}

	u.ID = lastID

	return &u, nil
}
