package redis

import (
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis(addr string) *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if _, err := RedisClient.Ping(RedisClient.Context()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil
	}
	log.Println("Connected to Redis successfully")
	return RedisClient
}
