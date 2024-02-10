package users_service

import (
	"pkg/service/pkg/data/request"
	"pkg/service/pkg/data/response"
)

type UsersService interface {
	SaveUserAction(req request.CreateUserActivityRequest) error
	GetUserActivities(req request.GetUserActivitiesRequest) (*response.GetUserActivitiesResponse, error)
}
