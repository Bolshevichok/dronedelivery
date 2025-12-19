package bootstrap

import (
	"fmt"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
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

	// Kafka producer для missions.created
	createdWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.MissionsCreatedTopic,
		Balancer: &kafka.LeastBytes{},
	}

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})

	return &missionService.Dependencies{
		Storage:     storage,
		KafkaWriter: createdWriter,
		RedisClient: redisClient,
	}, nil
}
