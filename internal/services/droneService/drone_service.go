package droneService

import (
	"context"
	"fmt"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/models"
	"github.com/segmentio/kafka-go"
)

type DroneStorage interface {
	GetMissionsByIDs(ctx context.Context, IDs []uint64) ([]*models.Mission, error)
	GetAvailableDrones(ctx context.Context, launchBaseID uint64) ([]*models.Drone, error)
	GetLaunchBasesByIDs(ctx context.Context, IDs []uint64) ([]*models.LaunchBase, error)
	UpdateMissionStatus(ctx context.Context, missionID uint64, status string) error
}

type DroneService struct {
	droneStorage    DroneStorage
	lifecycleWriter *kafka.Writer
	telemetryWriter *kafka.Writer
}

func NewDroneService(droneStorage DroneStorage, cfg *config.Config) *DroneService {
	lifecycleWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.DroneLifecycleTopic,
		Balancer: &kafka.LeastBytes{},
	}
	telemetryWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.DroneTelemetryTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &DroneService{
		droneStorage:    droneStorage,
		lifecycleWriter: lifecycleWriter,
		telemetryWriter: telemetryWriter,
	}
}
