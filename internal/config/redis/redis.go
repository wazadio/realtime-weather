package redis

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wazadio/realtime-weather/pkg/logger"
)

func NewRedisClient(ctx context.Context) *redis.Client {
	redisContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})

	_, err := redisClient.Ping(redisContext).Result()
	if err != nil {
		log.Fatalf("Error connecting to redis : %s", err.Error())
	}

	logger.Print(ctx, logger.INFO, "Redis connected")

	return redisClient
}
