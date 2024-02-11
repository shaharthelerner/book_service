package request

type UpdateBookTitleRequest struct {
	Id       string `binding:"required"`
	Title    string `json:"title" binding:"required"`
	Username string `json:"username" binding:"required"`
}
