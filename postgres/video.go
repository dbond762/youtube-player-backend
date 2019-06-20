package postgres

import (
	"database/sql"
	"github.com/lib/pq"
	"log"

	player "github.com/dbond762/youtube-player-backend"
	_ "github.com/lib/pq"
)

type VideoService struct {
	DB *sql.DB
}

func (s *VideoService) Like(u *player.User, v *player.Video) error {
	var lastID int64
	err := s.DB.QueryRow(
		`INSERT INTO "video" ("id", "title", "pub_date", "description", "thumbnail") VALUES ($1, $2, $3, $4, $5)`,
		v.ID, v.Title, v.PubDate, v.Description, v.Thumbnail).Scan(&lastID)
	if pgerr, ok := err.(*pq.Error); ok {
		if err != nil || pgerr.Code.Name() != "unique_violation" {
			log.Printf("Postgres: Error on scan row: %s", err)
			return err
		}
	}

	err = s.DB.QueryRow(`INSERT INTO "user_likes" ("id_user", "id_video") VALUES ($1, $2)`, u.ID, v.ID).Scan(&lastID)
	if err != nil {
		log.Printf("Postgres: Error on scan row: %s", err)
		return err
	}

	return nil
}

func (s *VideoService) Dislike(u *player.User, v *player.Video) error {
	_, err := s.DB.Exec(`DELETE FROM "user_likes" WHERE "id_user" = $1 AND "id_video"=$2`, u.ID, v.ID)
	if err != nil {
		log.Printf("Postgres: Error on delete: %s", err)
		return err
	}

	return nil
}

func (s *VideoService) IsLiked(u *player.User, v *player.Video) (bool, error) {
	row := s.DB.QueryRow(`SELECT "id" FROM "user_likes" WHERE "id_user" = $1 AND "id_video" = $2`, u.ID, v.ID)
	if err := row.Scan(&u.ID, &u.Login, &u.Password); err != nil {
		log.Printf("Postgres: Error on scan row: %s", err)
		return false, err
	}

	return true, nil
}
