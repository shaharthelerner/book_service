package request

type GetUserActivitiesRequest struct {
	Username string `form:"username" binding:"required"`
}
