package bootstrap

import (
	"context"
	"fmt"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/api/mission_api"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/mission_created_consumer"
	"github.com/Bolshevichok/dronedelivery/internal/services/droneService"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
	"github.com/Bolshevichok/dronedelivery/internal/services/processors/mission_processor"
	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
)

type MissionComponents struct {
	MissionService         missionService.MissionService
	MissionProcessor       mission_processor.MissionProcessor
	MissionCreatedConsumer mission_created_consumer.MissionCreatedConsumer
	MissionAPI             *mission_api.MissionAPI
}

func InitMissionComponents(cfg *config.Config) (*MissionComponents, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		fmt.Sprintf("%d", cfg.Database.Port),
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	storage, err := pgstorage.NewPGStorge(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to init storage: %w", err)
	}

	// Initialize services
	missionSvc := missionService.NewMissionService(context.Background(), storage, cfg)

	droneSvc := droneService.NewDroneService(storage, cfg)

	// Initialize processor
	processor := mission_processor.NewMissionProcessor(missionSvc, droneSvc)

	// Initialize consumer
	consumer := mission_created_consumer.NewMissionCreatedConsumer(processor, []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}, cfg.Kafka.MissionsCreatedTopic)

	// Initialize API
	missionAPI := mission_api.NewMissionAPI(missionSvc)

	return &MissionComponents{
		MissionService:         missionSvc,
		MissionProcessor:       processor,
		MissionCreatedConsumer: consumer,
		MissionAPI:             missionAPI,
	}, nil
}
