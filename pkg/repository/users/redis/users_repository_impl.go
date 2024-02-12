package users_repository

import (
	"context"
	"log"
	"pkg/service/pkg/models"
	"pkg/service/pkg/repository/users"
)

type RedisUsersRepository struct {
	ActivityActions int64
}

func NewRedisUsersRepository(activityActions int) users_repository.UsersRepository {
	return &RedisUsersRepository{
		ActivityActions: int64(activityActions),
	}
}

func (r *RedisUsersRepository) SaveAction(ua models.UserAction) error {
	client, err := NewRedisClient()
	if err != nil {
		return err
	}
	defer client.Close()

	key := createUsernameKey(ua.Username)

	// Push the action onto the left side of the list
	if err = client.LPush(context.Background(), key, ua.Action).Err(); err != nil {
		return err
	}
	// Trim the list to keep only the last r.ActivityActions elements
	if err = client.LTrim(context.Background(), key, 0, r.ActivityActions-1).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisUsersRepository) GetActivity(username string) (*models.UserActivity, error) {
	client, err := NewRedisClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	key := createUsernameKey(username)
	actions, err := client.LRange(context.Background(), key, 0, r.ActivityActions-1).Result()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &models.UserActivity{Actions: actions}, nil
}
