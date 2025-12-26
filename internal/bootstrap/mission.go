package bootstrap

import (
	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/api/mission_api"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/mission_created_consumer"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
	"github.com/Bolshevichok/dronedelivery/internal/services/processors/mission_processor"
)

type MissionComponents struct {
	MissionService         *missionService.MissionService
	MissionProcessor       mission_processor.MissionProcessor
	MissionCreatedConsumer mission_created_consumer.MissionCreatedConsumer
	MissionAPI             *mission_api.MissionAPI
}

func InitMissionComponents(cfg *config.Config) (*MissionComponents, error) {
	storage := InitPGStorage(cfg)

	// Initialize services
	missionSvc := InitMissionService(storage, cfg)
	droneSvc := InitDroneService(storage, cfg)

	// Initialize processor
	processor := InitMissionProcessor(droneSvc)

	// Initialize consumer
	consumer := InitMissionCreatedConsumer(cfg, processor)

	// Initialize API
	missionAPI := InitMissionAPI(missionSvc)

	return &MissionComponents{
		MissionService:         missionSvc,
		MissionProcessor:       processor,
		MissionCreatedConsumer: consumer,
		MissionAPI:             missionAPI,
	}, nil
}
