package postgres

import (
	"database/sql"
	"log"

	player "github.com/dbond762/youtube-player-backend"
	_ "github.com/lib/pq"
)

type VideoListService struct {
	DB     *sql.DB
	Finder player.VideoFinder
}

func (s *VideoListService) Search(query string) (*player.VideoList, error) {
	return s.Finder.Search(query)
}

func (s *VideoListService) SearchByUser(query string, user *player.User) (*player.VideoList, error) {
	list, err := s.Finder.Search(query)
	if err != nil {
		log.Printf("Postgres: error on search video: %s", err)
		return nil, err
	}

	// TODO: Add likes

	return list, nil
}

func (s *VideoListService) Likes(user *player.User) (*player.VideoList, error) {
	rows, err := s.DB.Query(`SELECT "video"."id", "video"."title", "video"."pub_date", "video"."description", "video"."thumbnail" FROM "user_likes" LEFT JOIN "video" ON ("user_likes"."id_video" = "video"."id") WHERE "user_likes"."id_user" = $1`, user.ID)
	if err != nil {
		log.Printf("Postgres: error on query: %s", err)
		return nil, err
	}

	const defaultCapacity = 25
	list := make(player.VideoList, 0, defaultCapacity)

	for rows.Next() {
		var v player.Video
		err := rows.Scan(&v.ID, &v.Title, &v.PubDate, &v.Description, &v.Thumbnail)
		if err != nil {
			log.Printf("Postgres: error on scan row: %s", err)
			return nil, err
		}
		list = append(list, player.VideoItem{Video: v})
	}
	rows.Close()

	return &list, nil
}
