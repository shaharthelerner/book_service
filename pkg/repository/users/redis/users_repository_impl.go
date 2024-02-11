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

func NewRedisUsersRepository(activityActions int64) users_repository.UsersRepository {
	return &RedisUsersRepository{
		ActivityActions: activityActions,
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

//func NewRedisClient() (*redis.Client, error) {
//	addr := os.Getenv("REDIS_ADDR")
//	if addr == "" {
//		addr = consts.DefaultRedisAddress
//	}
//	options := &redis.Options{
//		Addr:     addr,
//		Password: "",
//		DB:       0,
//	}
//	client := redis.NewClient(options)
//	if _, err := client.Ping(context.Background()).Result(); err != nil {
//		return nil, err
//	}
//
//	return client, nil
//}
//
//func createUsernameKey(username string) string {
//	return fmt.Sprintf(consts.UserActivitiesRedisKey, username)
//}