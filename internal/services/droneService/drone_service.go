package droneService

import (
	"fmt"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/segmentio/kafka-go"
)

type DroneService struct {
	lifecycleWriter *kafka.Writer
	telemetryWriter *kafka.Writer
}

func NewDroneService(cfg *config.Config) *DroneService {
	lifecycleWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.MissionsLifecycleTopic,
		Balancer: &kafka.LeastBytes{},
	}
	telemetryWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.DroneTelemetryTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &DroneService{
		lifecycleWriter: lifecycleWriter,
		telemetryWriter: telemetryWriter,
	}
}
