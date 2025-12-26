package main

import (
	"context"
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

	missionLifecycleProcessor := bootstrap.InitMissionLifecycleProcessor(missionService)
	missionLifecycleConsumer := bootstrap.InitMissionLifecycleConsumer(cfg, missionLifecycleProcessor)
	go missionLifecycleConsumer.Consume(context.Background())

	bootstrap.AppRun(missionApi)
}
