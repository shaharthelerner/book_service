package users_repository

import (
	"errors"
	"fmt"
	"log"
	"pkg/service/pkg/models"
	"pkg/service/pkg/repository/users"
)

type UsersRepositoryRedis struct {
	ActivityActions int64
}

func NewUsersRepositoryRedisImpl(activityActions int) users_repository.UsersRepository {
	return &UsersRepositoryRedis{
		ActivityActions: int64(activityActions),
	}
}

func (r *UsersRepositoryRedis) SaveAction(ua models.UserAction) error {
	client, err := r.newRedisClient()
	if err != nil {
		log.Printf("error creating redis client: %s", err)
		return err
	}
	defer client.Close()

	key := r.createUsernameKey(ua.Username)

	err = r.pushValueToKey(client, key, ua.Action)
	if err != nil {
		log.Printf("error saving action for user %s: %s", ua.Username, err)
		return err
	}

	err = r.trimKeyValues(client, key)
	if err != nil {
		log.Printf("error trimming action for user %s: %s", ua.Username, err)
	}

	return nil
}

func (r *UsersRepositoryRedis) GetActivity(username string) (*models.UserActivity, error) {
	client, err := r.newRedisClient()
	if err != nil {
		log.Printf("error creating redis client: %s", err)
		return nil, err
	}
	defer client.Close()

	key := r.createUsernameKey(username)
	actions, err := r.getRangeForKey(client, key)
	if err != nil {
		log.Printf("error getting activity for user %s: %s", username, err)
		return nil, errors.New(fmt.Sprintf("error getting activity for user %s", username))
	}

	return &models.UserActivity{Actions: actions}, nil
}
