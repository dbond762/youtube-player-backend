package http

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

/*
type ApiVideoResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				High struct {
					URL    string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
		Player struct {
			EmbedHTML string `json:"embedHtml"`
		} `json:"player"`
	} `json:"items"`
}
*/

type SearchResponse []struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pub_date"`
	Thumbnail   string    `json:"thumbnail"`
}

func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")

	list, err := h.VideoListService.Search(query)
	if err != nil {
		log.Printf("HTTP: Error on search video: %s", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	searchResp := make(SearchResponse, len(*list))
	for i, video := range *list {
		searchResp[i].ID = video.ID
		searchResp[i].Title = video.Title
		searchResp[i].Description = video.Description
		searchResp[i].PubDate = video.PubDate
		searchResp[i].Thumbnail = video.Thumbnail
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&searchResp); err != nil {
		log.Printf("HTTP: Error on encode search response: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
