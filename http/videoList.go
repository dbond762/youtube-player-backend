package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	player "github.com/dbond762/youtube-player-backend"
)

const ApiKey = "AIzaSyBVJgyC-x6CsM-hPCYY10VfOnGOKksDK8U"

type ApiSearchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

type VideoFinder struct{}

func (s *VideoFinder) Search(query string) (*player.VideoList, error) {

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&maxResults=25&q=%s&key=%s", query, ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP: error on get request: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	apiResp := new(ApiSearchResponse)
	if err := decoder.Decode(apiResp); err != nil {
		log.Printf("HTTP: error on decoding: %s", err)
		return nil, err
	}

	list := make(player.VideoList, len(apiResp.Items))
	for i, video := range apiResp.Items {
		v := player.Video{
			ID:          video.ID.VideoID,
			Title:       video.Snippet.Title,
			PubDate:     video.Snippet.PublishedAt,
			Description: video.Snippet.Description,
			Thumbnail:   video.Snippet.Thumbnails.High.URL,
		}
		list[i].Video = v
	}

	return &list, nil
}
