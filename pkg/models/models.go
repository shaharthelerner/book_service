package models

type Book struct {
	Id             string  `json:"id"`
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	Price          float64 `json:"price"`
	EbookAvailable bool    `json:"ebook_available"`
	PublishDate    string  `json:"publish_date"`
}

type BookFilters struct {
	Title      string
	AuthorName string
	MinPrice   float64
	MaxPrice   float64
}

type BookSource struct {
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	Price          float64 `json:"price"`
	EbookAvailable bool    `json:"ebook_available"`
	PublishDate    string  `json:"publish_date"`
}

type UserAction struct {
	Username string
	Action   string
}

type UserActivity struct {
	Actions []string
}

type StoreInventory struct {
	TotalBooks    int
	UniqueAuthors int
}
