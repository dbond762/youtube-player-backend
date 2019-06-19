package main

import (
	"database/sql"
	"log"

	"github.com/dbond762/youtube-player-backend/http"
	"github.com/dbond762/youtube-player-backend/postgres"
	"github.com/dbond762/youtube-player-backend/redis"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=youtube_player dbname=youtube_player password='123456' sslmode=disable")
	if err != nil {
		log.Fatalf("db err: %s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("db err: %s", err)
	}
	defer db.Close()

	us := &postgres.UserService{DB: db}
	session := redis.NewUserService("redis://youtube_player:@localhost:6379/0", us)
	vls := &postgres.VideoListService{
		DB:     db,
		Finder: redis.NewVideoSearcher("redis://youtube_player:@localhost:6379/1", new(http.VideoFinder)),
	}
	vs := &postgres.VideoService{DB: db}

	h := &http.Handler{
		UserService:      us,
		UserSession:      session,
		VideoListService: vls,
		VideoService:     vs,
	}

	http.Setup(h, 8080)
}
