package bootstrap

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
)

func InitMissionService(storage *pgstorage.PGstorage, cfg *config.Config) *missionService.MissionService {
	return missionService.NewMissionService(context.Background(), storage, cfg)
}
