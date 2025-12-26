package redis_consumer

import (
	"context"
	"log/slog"
	"time"
)

func (c *RedisConsumerImpl) Consume() {
	slog.Info("Starting Redis consumer on list", "list", c.channel)
	for {
		result, err := c.redisClient.BRPop(context.Background(), 0*time.Second, c.channel).Result()
		if err != nil {
			slog.Error("Error reading from Redis list", "error", err.Error())
			continue
		}
		if len(result) > 1 {
			message := result[1] // BRPop returns [key, value]
			slog.Info("Received message from Redis", "list", c.channel, "payload", message)
		}
	}
}
