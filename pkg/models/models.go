package models

type Book struct {
	Id             string  `json:"id"`
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	Price          float64 `json:"price"`
	EbookAvailable bool    `json:"ebook_available"`
	PublishDate    string  `json:"publish_date"`
	Username       string  `json:"username"`
}

type Store struct {
	TotalBooks int `json:"total_books"`
	Authors    int `json:"authors"`
}

type UserActivity struct {
	Method string
	Route  string
}
