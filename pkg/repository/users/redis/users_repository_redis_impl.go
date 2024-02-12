package users_repository

import (
	"context"
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

	// Push the action onto the left side of the list
	if err = client.LPush(context.Background(), key, ua.Action).Err(); err != nil {
		log.Printf("error saving action for user %s: %s", ua.Username, err)
		return err
	}

	// Trim the list to keep only the last r.ActivityActions elements
	if err = client.LTrim(context.Background(), key, 0, r.ActivityActions-1).Err(); err != nil {
		log.Printf("error trimming activity for user %s: %s", ua.Username, err)
		return err
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
	actions, err := client.LRange(context.Background(), key, 0, r.ActivityActions-1).Result()
	if err != nil {
		log.Printf("error getting activity for user %s: %s", username, err)
		return nil, errors.New(fmt.Sprintf("error getting activity for user %s", username))
	}

	return &models.UserActivity{Actions: actions}, nil
}
