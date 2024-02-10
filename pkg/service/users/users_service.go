package users_service

import (
	"pkg/service/pkg/data/request"
	"pkg/service/pkg/data/response"
)

type UsersService interface {
	CreateUserActivity(req request.CreateUserActivityRequest) error
	GetUserActivities(req request.GetUserActivitiesRequest) (*response.GetUserActivitiesResponse, error)
}
