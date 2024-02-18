package users_handler

import (
	"pkg/service/pkg/interfaces"
	"pkg/service/pkg/models"
	"pkg/service/pkg/models/request"
	"pkg/service/pkg/models/response"
)

var _ interfaces.UsersHandler = &UsersHandler{}

type UsersHandler struct {
	usersRepository interfaces.UsersRepository
}

func NewUsersHandler(usersRepository interfaces.UsersRepository) interfaces.UsersHandler {
	return &UsersHandler{
		usersRepository: usersRepository,
	}
}

func (u *UsersHandler) SaveUserAction(req request.CreateUserAction) error {
	userAction := models.UserAction{
		Username: req.Username,
		Action:   req.Method + " " + req.Route,
	}

	err := u.usersRepository.SaveAction(userAction)
	if err != nil {
		return err
	}

	return nil
}

func (u *UsersHandler) GetUserActivity(username string) (*response.GetUserActivity, error) {
	activity, err := u.usersRepository.GetActivity(username)
	if err != nil {
		return nil, err
	}

	return &response.GetUserActivity{Actions: activity.Actions}, nil
}
