package bootstrap

import (
	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/services/telemetryService"
)

func InitTelemetryService(cfg *config.Config) telemetryService.TelemetryService {
	return telemetryService.NewTelemetryService(cfg)
}
