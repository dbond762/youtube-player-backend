package youtube_player_backend

import (
	"time"
)

type SearchHistory struct {
	ID    int64
	User  *User
	Query string
	Time  time.Time
}

type SearchHistoryService interface {
	SearchHistory(u *User) (*SearchHistory, error)
}
