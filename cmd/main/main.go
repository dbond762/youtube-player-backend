package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/dbond762/youtube-player-backend/http"
	"github.com/dbond762/youtube-player-backend/postgres"
	"log"
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

	var h http.Handler
	h.UserService = us

	http.Setup(h, 8080)
}
