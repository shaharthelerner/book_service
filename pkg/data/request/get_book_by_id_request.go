package request

type GetBookByIdRequest struct {
	Id       string `form:"id" binding:"required"`
	Username string `form:"username" binding:"required"`
}
