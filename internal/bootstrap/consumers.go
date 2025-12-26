package bootstrap

import (
	"fmt"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/mission_created_consumer"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/telemetry_consumer"
	missionprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/mission_processor"
	telemetryprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/telemetry_processor"
)

func InitMissionCreatedConsumer(cfg *config.Config, missionProcessor *missionprocessor.MissionProcessorImpl) mission_created_consumer.MissionCreatedConsumer {
	kafkaBrokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return mission_created_consumer.NewMissionCreatedConsumer(missionProcessor, kafkaBrokers, cfg.Kafka.MissionsCreatedTopic)
}

func InitTelemetryConsumer(cfg *config.Config, telemetryProcessor *telemetryprocessor.TelemetryProcessorImpl) telemetry_consumer.TelemetryConsumer {
	kafkaBrokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return telemetry_consumer.NewTelemetryConsumer(telemetryProcessor, kafkaBrokers, cfg.Kafka.DroneTelemetryTopic)
}
