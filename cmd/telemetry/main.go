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

	telemetryService := bootstrap.InitTelemetryService(cfg)
	telemetryProcessor := bootstrap.InitTelemetryProcessor(telemetryService)
	telemetryConsumer := bootstrap.InitTelemetryConsumer(cfg, telemetryProcessor)

	bootstrap.AppRunConsumers(telemetryConsumer)
}
