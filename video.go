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

type VideoService interface {
	Video(id string) (*Video, error)
	Like(u *User, v *Video) error
	Dislike(u *User, v *Video) error
	IsLiked(u *User, v *Video) (bool, error)
}
