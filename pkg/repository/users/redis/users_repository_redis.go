package redis

import (
	"errors"
	"fmt"
	"log"
	"pkg/service/pkg/interfaces"
	"pkg/service/pkg/models"
)

var _ interfaces.UsersRepository = &UsersRepositoryRedis{}

type UsersRepositoryRedis struct {
	activityActions int64
}

func NewUsersRepositoryRedis(activityActions int) interfaces.UsersRepository {
	return &UsersRepositoryRedis{
		activityActions: int64(activityActions),
	}
}

func (r *UsersRepositoryRedis) SaveAction(ua models.UserAction) error {
	client, err := newRedisClient()
	if err != nil {
		log.Printf("error creating redis client: %s", err)
		return err
	}
	defer client.Close()

	key := createUsernameKey(ua.Username)

	err = pushValueToKey(client, key, ua.Action)
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
	client, err := newRedisClient()
	if err != nil {
		log.Printf("error creating redis client: %s", err)
		return nil, err
	}
	defer client.Close()

	key := createUsernameKey(username)
	actions, err := r.getRangeForKey(client, key)
	if err != nil {
		log.Printf("error getting activity for user %s: %s", username, err)
		return nil, errors.New(fmt.Sprintf("error getting activity for user %s", username))
	}

	return &models.UserActivity{Actions: actions}, nil
}
