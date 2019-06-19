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
		return nil, err
	}

	rows, err := s.DB.Query(`SELECT "video"."id" FROM "user_likes" LEFT JOIN "video" ON ("user_likes"."id_video" = "video"."id") WHERE "user_likes"."id_user" = $1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var v player.Video
		err := rows.Scan(&v.ID)
		if err != nil {
			return nil, err
		}
		u.Likes = append(u.Likes, v)
	}
	rows.Close()

	return &u, nil
}

func (s *UserService) UserByLogin(login string) (*player.User, error) {
	var u player.User
	row := s.DB.QueryRow(`SELECT "id", "login", "password" FROM "user" WHERE "login" = $1`, login)
	if err := row.Scan(&u.ID, &u.Login, &u.Password); err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(`SELECT "video"."id" FROM "user_likes" LEFT JOIN "video" ON ("user_likes"."id_video" = "video"."id") WHERE "user_likes"."id_user" = $1`, u.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var v player.Video
		err := rows.Scan(&v.ID)
		if err != nil {
			return nil, err
		}
		u.Likes = append(u.Likes, v)
	}
	rows.Close()

	return &u, nil
}

func (s *UserService) CreateUser(u player.User) (*player.User, error) {
	var lastID int64
	hash, err := player.HashPassword([]byte(u.Password))
	if err != nil {
		return nil, err
	}

	err = s.DB.QueryRow(`INSERT INTO "user" ("login", "password") VALUES ($1, $2)  RETURNING "id"`, u.Login, hash).Scan(&lastID)
	if err != nil {
		log.Print(lastID)
		return nil, err
	}

	u.ID = lastID

	rows, err := s.DB.Query(`SELECT "video"."id" FROM "user_likes" LEFT JOIN "video" ON ("user_likes"."id_video" = "video"."id") WHERE "user_likes"."id_user" = $1`, u.ID)
	if err != nil {
		log.Printf("rows")
		return nil, err
	}

	for rows.Next() {
		var v player.Video
		err := rows.Scan(&v.ID)
		if err != nil {
			log.Printf("rows scan")
			return nil, err
		}
		u.Likes = append(u.Likes, v)
	}
	rows.Close()

	return &u, nil
}
