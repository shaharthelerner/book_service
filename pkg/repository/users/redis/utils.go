package users_repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"pkg/service/pkg/consts"
)

func (r *UsersRepositoryRedis) newRedisClient() (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = consts.DefaultRedisAddress
	}
	options := &redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	}

	client := redis.NewClient(options)
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return client, nil
}

func (r *UsersRepositoryRedis) createUsernameKey(username string) string {
	return fmt.Sprintf(consts.UserActivitiesRedisKey, username)
}
