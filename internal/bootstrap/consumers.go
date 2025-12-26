package bootstrap

import (
	"fmt"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/mission_created_consumer"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/mission_lifecycle_consumer"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/redis_consumer"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/telemetry_consumer"
	missionlifecycleprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/mission_lifecycle_processor"
	missionprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/mission_processor"
	telemetryprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/telemetry_processor"
	"github.com/go-redis/redis/v8"
)

func InitMissionCreatedConsumer(cfg *config.Config, missionProcessor *missionprocessor.MissionProcessorImpl) mission_created_consumer.MissionCreatedConsumer {
	kafkaBrokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return mission_created_consumer.NewMissionCreatedConsumer(missionProcessor, kafkaBrokers, cfg.Kafka.MissionsCreatedTopic)
}

func InitTelemetryConsumer(cfg *config.Config, telemetryProcessor *telemetryprocessor.TelemetryProcessorImpl) telemetry_consumer.TelemetryConsumer {
	kafkaBrokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return telemetry_consumer.NewTelemetryConsumer(telemetryProcessor, kafkaBrokers, cfg.Kafka.DroneTelemetryTopic)
}

func InitMissionLifecycleConsumer(cfg *config.Config, processor *missionlifecycleprocessor.MissionLifecycleProcessorImpl) mission_lifecycle_consumer.MissionLifecycleConsumer {
	kafkaBrokers := []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)}
	return mission_lifecycle_consumer.NewMissionLifecycleConsumer(processor, kafkaBrokers, cfg.Kafka.MissionsLifecycleTopic)
}

func InitRedisConsumer(redisClient *redis.Client, channel string) redis_consumer.RedisConsumer {
	return redis_consumer.NewRedisConsumer(redisClient, channel)
}
