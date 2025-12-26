package redis_consumer

import (
	"github.com/go-redis/redis/v8"
)

type RedisConsumer interface {
	Consume()
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
