package request

type UpdateBookTitleRequest struct {
	Id       string `json:"id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Username string `json:"username" binding:"required"`
}
