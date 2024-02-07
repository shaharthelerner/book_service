package books

type store struct {
	Books   int `json:"books"`
	Authors int `json:"authors"`
}

type GetBooksRequest struct {
	Title      string  `form:"title"`
	AuthorName string  `form:"author_name"`
	MinPrice   float64 `form:"min_price" validate:"gte=0"`
	MaxPrice   float64 `form:"max_price" validate:"gte=0"`
	Username   string  `form:"username" binding:"required"`
}

type GetBookRequest struct {
	Id       string `form:"id" binding:"required"`
	Username string `form:"username" binding:"required"`
}

type GetBookResponse struct {
	Found  bool     `json:"found"`
	Source BookRead `json:"_source"`
}

type CreateBookRequest struct {
	Title          string  `json:"title" binding:"required"`
	AuthorName     string  `json:"author_name" binding:"required"`
	Price          float64 `json:"price" binding:"required"`
	EbookAvailable bool    `json:"ebook_available" binding:"required"`
	PublishDate    string  `json:"publish_date" binding:"required"`
	Username       string  `json:"username" binding:"required"`
}

type CreateBookResponse struct {
	Result string                 `json:"result,omitempty"`
	Error  map[string]interface{} `json:"error,omitempty"`
	Status int                    `json:"status,omitempty"`
}

type CreateBookObject struct {
	Id             string  `json:"id"`
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	Price          float64 `json:"price"`
	EbookAvailable bool    `json:"ebook_available"`
	PublishDate    string  `json:"publish_date"`
	Username       string  `json:"username"`
}

type UpdateBookRequest struct {
	Id       string `json:"id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type UpdateBookResponse struct {
	Result string                 `json:"result,omitempty"`
	Error  map[string]interface{} `json:"error,omitempty"`
	Status int                    `json:"status,omitempty"`
}

type DeleteBookRequest struct {
	Id       string `form:"id" binding:"required"`
	Username string `form:"username" binding:"required"`
}

type DeleteBookResponse struct {
	Result string `json:"result"`
}

type GetStoreInventoryRequest struct {
	Username string `form:"username" binding:"required"`
}

type GetBooksResponse struct {
	Hits Hits `json:"hits"`
}

type UserActivityRequest struct {
	Username string `form:"username" binding:"required"`
}

type Hits struct {
	Hits []BookHit `json:"hits"`
}

type BookHit struct {
	Source BookRead `json:"_source"`
}

type BookRead struct {
	Title          string  `json:"title"`
	AuthorName     string  `json:"author_name"`
	Price          float64 `json:"price"`
	EbookAvailable bool    `json:"ebook_available"`
	PublishDate    string  `json:"publish_date"`
}

type UserActivity struct {
	Method string
	Route  string
}
