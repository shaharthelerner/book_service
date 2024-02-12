package request

type GetStoreInventoryRequest struct {
	Username string `form:"username" binding:"required"`
}
