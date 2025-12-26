package bootstrap

import (
	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/services/droneService"
	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
)

func InitDroneService(storage *pgstorage.PGstorage, cfg *config.Config) *droneService.DroneService {
	return droneService.NewDroneService(storage, cfg)
}
