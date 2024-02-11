package request

type GetBookByIdRequest struct {
	Id       string `binding:"required"`
	Username string `form:"username" binding:"required"`
}
