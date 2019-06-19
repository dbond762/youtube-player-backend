package youtube_player_backend

import (
	"time"
)

type Video struct {
	ID          string
	Title       string
	PubDate     time.Time
	Description string
	Thumbnail   string
	Player      string
}
