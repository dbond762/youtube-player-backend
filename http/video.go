package http

import (
	"encoding/json"
	"log"
	"net/http"

	player "github.com/dbond762/youtube-player-backend"
	"github.com/go-chi/chi"
)

func (h *Handler) like(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isLoggedIn, ok := ctx.Value("isLoggedIn").(bool)
	if !ok {
		log.Printf("HTTP: isLoggedIn not found in context")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !isLoggedIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	user, ok := ctx.Value("user").(*player.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	videoID := chi.URLParam(r, "videoID")

	video, err := h.VideoService.Video(videoID)
	if err != nil {
		log.Printf("HTTP: Not found video: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.VideoService.Like(user, video)
	if err != nil {
		log.Printf("HTTP: Coudn`t like video: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) dislike(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isLoggedIn, ok := ctx.Value("isLoggedIn").(bool)
	if !ok {
		log.Printf("HTTP: isLoggedIn not found in context")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !isLoggedIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	user, ok := ctx.Value("user").(*player.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	videoID := chi.URLParam(r, "videoID")

	video, err := h.VideoService.Video(videoID)
	if err != nil {
		log.Printf("HTTP: Not found video: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.VideoService.Dislike(user, video)
	if err != nil {
		log.Printf("HTTP: Coudn`t dislike video: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func (h *Handler) getLikes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isLoggedIn, ok := ctx.Value("isLoggedIn").(bool)
	if !ok {
		log.Printf("HTTP: isLoggedIn not found in context")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !isLoggedIn {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	user, ok := ctx.Value("user").(*player.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	list, err := h.VideoListService.Likes(user)
	if err != nil {
		log.Printf("HTTP: Error on get likes: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := make(SearchResponse, len(*list))
	for i, video := range *list {
		resp[i].ID = video.ID
		resp[i].Title = video.Title
		resp[i].Description = video.Description
		resp[i].PubDate = video.PubDate
		resp[i].Thumbnail = video.Thumbnail
		resp[i].Liked = video.Liked
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&resp); err != nil {
		log.Printf("HTTP: Error on encode response: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
