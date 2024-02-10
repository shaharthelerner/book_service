package users_database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

type UsersActivityRedis struct {
	client         *redis.Client
	actionsToCache int64
}

const Key = "book_service_exercise:users:activities:%s"

func NewUsersActivityRedis(maxUserActions int64) (*UsersActivityRedis, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	options := &redis.Options{
		Addr:     addr, // Replace with your Redis server address
		Password: "",   // No password by default
		DB:       0,    // Default DB
	}
	client := redis.NewClient(options)
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &UsersActivityRedis{
		client:         client,
		actionsToCache: maxUserActions,
	}, nil
}

func (r *UsersActivityRedis) CreateUserAction(username string, action string) error {
	key := fmt.Sprintf(Key, username)

	// Push the action onto the left side of the list
	err := r.client.LPush(context.Background(), key, action).Err()
	if err != nil {
		log.Fatal(err)
	}

	// Trim the list to keep only the last MaxActions elements
	err = r.client.LTrim(context.Background(), key, 0, r.actionsToCache-1).Err()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *UsersActivityRedis) GetUserActivity(username string) ([]string, error) {
	key := fmt.Sprintf(Key, username)

	activities, err := r.client.LRange(context.Background(), key, 0, r.actionsToCache-1).Result()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return activities, nil
}

func (r *UsersActivityRedis) Close() error {
	return r.client.Close()
}

/*
const User = "user3"

func SaveMockCache() {
	client, err := NewUsersActivityRedis(MaxActions)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func(client *UsersActivityRedis) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	// Record user actions
	if err := client.CreateUserAction(User, "GET /api/resource11"); err != nil {
		log.Fatal(err)
	}
	if err := client.CreateUserAction(User, "POST /api/resource21"); err != nil {
		log.Fatal(err)
	}
	if err := client.CreateUserAction(User, "PUT /api/resource31"); err != nil {
		log.Fatal(err)
	}
	if err := client.CreateUserAction(User, "PUT /api/resource41"); err != nil {
		log.Fatal(err)
	}
	if err := client.CreateUserAction(User, "PUT /api/resource51"); err != nil {
		log.Fatal(err)
	}
	if err := client.CreateUserAction(User, "PUT /api/resource61"); err != nil {
		log.Fatal(err)
	}

	actions, err := client.GetUserActivity(User)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Last actions for", User, ":", actions)
	}
}
*/
