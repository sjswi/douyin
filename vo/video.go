package vo

type Video struct {
	Author        *User  `json:"author"`
	ID            uint   `json:"id"`
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
	PlayURL       string `json:"play_url"`
	CoverURL      string `json:"cover_url"`
}
