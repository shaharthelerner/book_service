package users_repository

import (
	"log"
	database "pkg/service/pkg/database/users"
	"pkg/service/pkg/models"
)

type UsersRepositoryImpl struct {
	ActivityActions int64
}

func NewUsersRepositoryImpl(activityActions int64) UsersRepository {
	return &UsersRepositoryImpl{
		ActivityActions: activityActions,
	}
}

func (ur *UsersRepositoryImpl) SaveAction(ua models.UserAction) error {
	r, err := database.NewUsersRedis(ur.ActivityActions)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = r.SaveAction(ua.Username, ua.Action)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UsersRepositoryImpl) GetActivity(username string) (*models.UserActivity, error) {
	r, err := database.NewUsersRedis(ur.ActivityActions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	actions, err := r.GetUserActivity(username)
	if err != nil {
		return nil, err
	}

	return &models.UserActivity{Actions: actions}, nil
}
