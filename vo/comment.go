package vo

type Comment struct {
	ID         uint   `json:"id"`
	User       *User  `json:"user"`
	Content    string `json:"content"`
	CreateDate string `json:"create_date"`
}
