package request

type UpdateBookTitle struct {
	Title string `json:"title" binding:"required"`
}
