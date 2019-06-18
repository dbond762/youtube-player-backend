package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"time"
)

const ApiKey = "AIzaSyBVJgyC-x6CsM-hPCYY10VfOnGOKksDK8U"

type ApiResponse struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
			} `json:"thumbnails"`
			ChannelTitle         string `json:"channelTitle"`
			LiveBroadcastContent string `json:"liveBroadcastContent"`
		} `json:"snippet"`
	} `json:"items"`
}

type SearchResponse []struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	PubDate time.Time `json:"pub_date"`
	Thumbnails  struct {
		Default struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"default"`
		Medium struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"medium"`
		High struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"high"`
	} `json:"thumbnails"`
}

func search(w http.ResponseWriter, r *http.Request) {
	query := chi.URLParam(r, "query")

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&maxResults=20&q=%s&key=%s", query, ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		//
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	apiResp := new(ApiResponse)
	if err := decoder.Decode(apiResp); err != nil {
		//
	}

	searchResp := make(SearchResponse, len(apiResp.Items))
	for i, video := range apiResp.Items {
		searchResp[i].ID = video.ID.VideoID
		searchResp[i].Name = video.Snippet.Title
		searchResp[i].Description = video.Snippet.Description
		searchResp[i].PubDate = video.Snippet.PublishedAt
		searchResp[i].Thumbnails = video.Snippet.Thumbnails
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&searchResp); err != nil {
		//
	}
}
