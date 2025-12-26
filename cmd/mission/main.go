package main

import (
	"fmt"
	"os"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	storage := bootstrap.InitPGStorage(cfg)
	missionService := bootstrap.InitMissionService(storage, cfg)
	missionApi := bootstrap.InitMissionAPI(missionService)

	redisClient := bootstrap.InitRedis(cfg)
	redisConsumer := bootstrap.InitRedisConsumer(redisClient, "telemetry_queue")

	missionLifecycleProcessor := bootstrap.InitMissionLifecycleProcessor(missionService)
	missionLifecycleConsumer := bootstrap.InitMissionLifecycleConsumer(cfg, missionLifecycleProcessor)

	bootstrap.AppRun(missionApi, redisConsumer, missionLifecycleConsumer)
}
