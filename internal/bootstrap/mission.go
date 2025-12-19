package bootstrap

import (
	"fmt"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
)

// InitMissionService initializes dependencies for mission-service
func InitMissionService(cfg *config.Config) (*missionService.Dependencies, error) {
	// Build DB connection string
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		fmt.Sprintf("%d", cfg.Database.Port),
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// Initialize storage
	storage, err := pgstorage.NewPGStorge(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to init storage: %w", err)
	}

	// TODO: Initialize Kafka producer/consumer, Redis, etc.

	return &missionService.Dependencies{
		Storage: storage,
	}, nil
}
