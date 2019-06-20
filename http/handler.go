package http

import (
	"fmt"
	"log"
	"net/http"

	player "github.com/dbond762/youtube-player-backend"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	UserService      player.UserService
	UserSession      player.UserSession
	VideoListService player.VideoListService
	VideoService     player.VideoService
}

func Setup(h *Handler, port int) {
	r := chi.NewRouter()

	CORS := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(CORS.Handler)

	r.Use(h.AuthMiddleware)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/search/", h.search)
	r.Get("/search/{query}", h.search)

	r.Post("/signup", h.signUp)
	r.Post("/login", h.login)
	r.Get("/logout", h.logout)

	r.Get("/like/{videoID}", h.like)
	r.Get("/dislike/{videoID}", h.dislike)

	log.Printf("Server run on http://localhost:%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatal("HTTP: err on ListenAndServe: ", err)
	}
}
