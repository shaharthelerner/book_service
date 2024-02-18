package request

type GetBooks struct {
	Title      string   `form:"title"`
	AuthorName string   `form:"author_name"`
	MinPrice   *float64 `form:"min_price"`
	MaxPrice   *float64 `form:"max_price"`
}
