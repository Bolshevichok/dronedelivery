package bootstrap

import (
	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/services/droneService"
)

func InitDroneService(cfg *config.Config) *droneService.DroneService {
	return droneService.NewDroneService(cfg)
}
