package interfaces

import (
	"pkg/service/pkg/models/request"
	"pkg/service/pkg/models/response"
)

type UsersHandler interface {
	SaveUserAction(req request.CreateUserAction) error
	GetUserActivity(username string) (*response.GetUserActivity, error)
}
