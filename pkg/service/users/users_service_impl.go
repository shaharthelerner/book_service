package users_service

import (
	"log"
	"pkg/service/pkg/data/request"
	"pkg/service/pkg/data/response"
	database "pkg/service/pkg/database/users"
	repository "pkg/service/pkg/repository/users"
)

const MaxActions = 3

type UsersServiceImpl struct {
	UsersRepository repository.UsersRepository
}

func NewUsersServiceImpl(usersRepository repository.UsersRepository) UsersService {
	return &UsersServiceImpl{
		UsersRepository: usersRepository,
	}
}

func (us *UsersServiceImpl) CreateUserActivity(req request.CreateUserActivityRequest) error {
	r, err := database.NewUsersActivityRedis(MaxActions)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = r.CreateUserAction(req.Username, req.Activity.Method+" "+req.Activity.Route)
	if err != nil {
		log.Fatalf("Error saving user activity: %v", err)
		return err
	}
	return nil
}

func (us *UsersServiceImpl) GetUserActivities(req request.GetUserActivitiesRequest) (*response.GetUserActivitiesResponse, error) {
	r, err := database.NewUsersActivityRedis(MaxActions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	actions, err := r.GetUserActivity(req.Username)
	if err != nil {
		log.Fatalf("Error fetching user actions: %v", err)
		return nil, err
	}
	return &response.GetUserActivitiesResponse{Actions: actions}, nil
}
