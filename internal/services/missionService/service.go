package missionService

import (
	"context"
	"fmt"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
)

type MissionStorage interface {
	UpsertMissions(ctx context.Context, missions []*models.Mission) ([]*models.Mission, error)
	GetMissionsByIDs(ctx context.Context, IDs []uint64) ([]*models.Mission, error)
	UpdateMissionStatus(ctx context.Context, missionID uint64, status string) error
}

type MissionService struct {
	mission_api.UnimplementedMissionServiceServer
	missionStorage MissionStorage
	kafkaWriter    *kafka.Writer
	redisClient    *redis.Client
}

func NewMissionService(ctx context.Context, missionStorage MissionStorage, cfg *config.Config) *MissionService {
	createdWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.MissionsCreatedTopic,
		Balancer: &kafka.LeastBytes{},
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})

	return &MissionService{
		missionStorage: missionStorage,
		kafkaWriter:    createdWriter,
		redisClient:    redisClient,
	}
}

func (s *MissionService) GetDroneTelemetry(ctx context.Context, droneID string) (string, error) {
	key := fmt.Sprintf("drone:%s", droneID)
	val, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get telemetry for drone %s: %w", droneID, err)
	}
	return val, nil
}
