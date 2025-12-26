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
	droneService := bootstrap.InitDroneService(storage, cfg)
	missionProcessor := bootstrap.InitMissionProcessor(droneService)
	missionCreatedConsumer := bootstrap.InitMissionCreatedConsumer(cfg, missionProcessor)

	bootstrap.AppRunConsumer(missionCreatedConsumer)
}
