package redis_consumer

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisConsumer interface {
	Consume(ctx context.Context)
}

type RedisConsumerImpl struct {
	redisClient *redis.Client
	channel     string
}

func NewRedisConsumer(redisClient *redis.Client, channel string) *RedisConsumerImpl {
	return &RedisConsumerImpl{
		redisClient: redisClient,
		channel:     channel,
	}
}
