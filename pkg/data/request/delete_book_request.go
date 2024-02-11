package request

type DeleteBookRequest struct {
	Id       string `binding:"required"`
	Username string `form:"username" binding:"required"`
}
