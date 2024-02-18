package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"pkg/service/pkg/consts"
	"time"
)

func newRedisClient() (*redis.Client, error) {
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

func createUsernameKey(username string) string {
	return fmt.Sprintf(consts.UserActivityRedisKey, username)
}

func pushValueToKey(client *redis.Client, key string, value any) error {
	ctx, cancel := context.WithTimeout(context.Background(), consts.UsersRequestTimeout*time.Second)
	defer cancel()
	return client.LPush(ctx, key, value).Err()
}

func (r *UsersRepositoryRedis) trimKeyValues(client *redis.Client, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), consts.UsersRequestTimeout*time.Second)
	defer cancel()
	return client.LTrim(ctx, key, 0, r.activityActions-1).Err()
}

func (r *UsersRepositoryRedis) getRangeForKey(client *redis.Client, key string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), consts.UsersRequestTimeout*time.Second)
	defer cancel()
	return client.LRange(ctx, key, 0, r.activityActions-1).Result()
}
