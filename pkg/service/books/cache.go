package books

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

type RedisCacheBooks struct {
	MaxUserActions int64
}

const MaxActions = 3
const Key = "book_service_exercise:users:activities:%s"
const User = "user3"

func SaveMockCache() {
	r := newCache(MaxActions)
	client, err := connectToRedis()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Record user actions
	if err := r.SetUserActivity(client, User, "GET /api/resource11"); err != nil {
		log.Fatal(err)
	}
	if err := r.SetUserActivity(client, User, "POST /api/resource21"); err != nil {
		log.Fatal(err)
	}
	if err := r.SetUserActivity(client, User, "PUT /api/resource31"); err != nil {
		log.Fatal(err)
	}
	if err := r.SetUserActivity(client, User, "PUT /api/resource41"); err != nil {
		log.Fatal(err)
	}
	if err := r.SetUserActivity(client, User, "PUT /api/resource51"); err != nil {
		log.Fatal(err)
	}
	if err := r.SetUserActivity(client, User, "PUT /api/resource61"); err != nil {
		log.Fatal(err)
	}

	actions, err := r.GetUserActivity(client, User)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Last actions for", User, ":", actions)
	}

	// Close the Redis client
	if err := closeRedis(client); err != nil {
		log.Fatal(err)
	}
}

func newCache(maxUserActions int64) *RedisCacheBooks {
	return &RedisCacheBooks{
		MaxUserActions: maxUserActions,
	}
}

func (r *RedisCacheBooks) SetUserActivity(client *redis.Client, username string, action string) error {
	key := fmt.Sprintf(Key, username)

	// Push the action onto the left side of the list
	err := client.LPush(context.Background(), key, action).Err()
	if err != nil {
		log.Fatal(err)
	}

	// Trim the list to keep only the last MaxActions elements
	err = client.LTrim(context.Background(), key, 0, r.MaxUserActions-1).Err()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *RedisCacheBooks) GetUserActivity(client *redis.Client, username string) ([]string, error) {
	key := fmt.Sprintf(Key, username)

	actions, err := client.LRange(context.Background(), key, 0, r.MaxUserActions-1).Result()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return actions, nil
}

func connectToRedis() (client *redis.Client, err error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	options := &redis.Options{
		Addr:     addr, // Replace with your Redis server address
		Password: "",   // No password by default
		DB:       0,    // Default DB
	}

	// Create a new Redis client
	client = redis.NewClient(options)

	if _, err = client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return client, nil

}

func closeRedis(client *redis.Client) error {
	return client.Close()
}
