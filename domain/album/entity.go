package album

type Album struct {
	UserID int    `json:"userId"`
	ID     int64  `json:"id"`
	Title  string `json:"title"`
}
