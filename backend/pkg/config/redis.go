package config

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func CreateRedisClient() {
	// opt, err := redis.ParseURL("redis://localhost:6379/0")
	opt, err := redis.ParseURL(os.Getenv("REDISCLOUD_URL"))
	if err != nil {
		panic(err)
	}
	redis := redis.NewClient(opt)
	Redis = redis

}
