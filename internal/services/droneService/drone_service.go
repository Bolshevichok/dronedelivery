package droneService

import (
	"context"
	"fmt"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal"
	"github.com/segmentio/kafka-go"
)

// DroneService interface
type DroneService interface {
	ProcessMissionCreated(ctx context.Context, missionID uint64)
}

type DroneServiceImpl struct {
	storage         internal.Storage
	lifecycleWriter *kafka.Writer
	telemetryWriter *kafka.Writer
}

func NewDroneService(storage internal.Storage, cfg *config.Config) DroneService {
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

	return &DroneServiceImpl{
		storage:         storage,
		lifecycleWriter: lifecycleWriter,
		telemetryWriter: telemetryWriter,
	}
}
