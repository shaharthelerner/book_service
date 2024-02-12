package request

type GetUserActivitiesRequest struct {
	Username string `binding:"required"`
}
