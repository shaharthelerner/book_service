package request

type UserActivityRequest struct {
	Username string `form:"username" binding:"required"`
}
