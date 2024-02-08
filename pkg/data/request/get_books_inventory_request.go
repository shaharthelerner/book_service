package request

type GetBooksInventoryRequest struct {
	Username string `form:"username" binding:"required"`
}
