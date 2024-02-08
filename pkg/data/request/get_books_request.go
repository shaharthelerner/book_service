package request

type GetBooksRequest struct {
	Title      string  `form:"title"`
	AuthorName string  `form:"author_name"`
	MinPrice   float64 `form:"min_price" validate:"gte=0"`
	MaxPrice   float64 `form:"max_price" validate:"gte=0"`
	Username   string  `form:"username" binding:"required"`
}
