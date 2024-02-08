package request

type CreateBookRequest struct {
	Title          string  `json:"title" binding:"required"`
	AuthorName     string  `json:"author_name" binding:"required"`
	Price          float64 `json:"price" binding:"required"`
	EbookAvailable bool    `json:"ebook_available" binding:"required"`
	PublishDate    string  `json:"publish_date" binding:"required"`
	Username       string  `json:"username" binding:"required"`
}
