package youtube_player_backend

type VideoItem struct {
	Video
	Liked bool
}

type VideoList []VideoItem

type VideoFinder interface {
	Search(query string) (*VideoList, error)
}

type VideoListService interface {
	VideoFinder
	SearchByUser(query string, user *User) (*VideoList, error)
	Likes(user *User) (*VideoList, error)
}
