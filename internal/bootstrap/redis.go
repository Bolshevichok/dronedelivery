package bootstrap

import (
	"fmt"
	"log"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/go-redis/redis/v8"
)

func InitRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})
	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		log.Panic(fmt.Sprintf("ошибка подключения к Redis, %v", err))
		panic(err)
	}
	return rdb
}
