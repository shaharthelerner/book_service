package users_service

import (
	"pkg/service/pkg/data/request"
	"pkg/service/pkg/data/response"
	"pkg/service/pkg/models"
	repository "pkg/service/pkg/repository/users"
)

type UsersServiceImpl struct {
	UsersRepository repository.UsersRepository
}

func NewUsersServiceImpl(usersRepository repository.UsersRepository) UsersService {
	return &UsersServiceImpl{
		UsersRepository: usersRepository,
	}
}

func (us *UsersServiceImpl) SaveUserAction(req request.CreateUserActivityRequest) error {
	ua := models.UserAction{
		Username: req.Username,
		Action:   req.Method + " " + req.Route,
	}

	err := us.UsersRepository.SaveAction(ua)
	if err != nil {
		return err
	}

	return nil
}

func (us *UsersServiceImpl) GetUserActivities(req request.GetUserActivitiesRequest) (*response.GetUserActivitiesResponse, error) {
	activity, err := us.UsersRepository.GetActivity(req.Username)
	if err != nil {
		return nil, err
	}

	return &response.GetUserActivitiesResponse{Actions: activity.Actions}, nil
}
