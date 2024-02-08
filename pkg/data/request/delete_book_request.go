package request

type DeleteBookRequest struct {
	Id       string `form:"id" binding:"required"`
	Username string `form:"username" binding:"required"`
}
