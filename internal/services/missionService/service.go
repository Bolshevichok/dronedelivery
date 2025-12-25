package missionService

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal"
	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
	models_pb "github.com/Bolshevichok/dronedelivery/internal/pb/models"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MissionService interface
type MissionService interface {
	UpsertMissions(ctx context.Context, req *mission_api.UpsertMissionsRequest) (*mission_api.UpsertMissionsResponse, error)
	GetMission(ctx context.Context, req *mission_api.GetMissionRequest) (*mission_api.GetMissionResponse, error)
}

type Dependencies struct {
	Storage     internal.Storage
	KafkaWriter *kafka.Writer
	RedisClient *redis.Client
}

type MissionServiceImpl struct {
	mission_api.UnimplementedMissionServiceServer
	deps *Dependencies
}

func NewMissionService(ctx context.Context, storage internal.Storage, cfg *config.Config) MissionService {
	createdWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.MissionsCreatedTopic,
		Balancer: &kafka.LeastBytes{},
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})

	deps := &Dependencies{
		Storage:     storage,
		KafkaWriter: createdWriter,
		RedisClient: redisClient,
	}

	return &MissionServiceImpl{deps: deps}
}

func (s *MissionServiceImpl) UpsertMissions(ctx context.Context, req *mission_api.UpsertMissionsRequest) (*mission_api.UpsertMissionsResponse, error) {
	for _, missionPb := range req.Missions {
		mission := &models.Mission{
			OperatorID:     missionPb.OpId,
			LaunchBaseID:   missionPb.BaseId,
			Status:         missionPb.Status,
			DestinationLat: missionPb.Lat,
			DestinationLon: missionPb.Lon,
			DestinationAlt: missionPb.Alt,
			PayloadKg:      missionPb.Payload,
			CreatedAt:      time.Now(),
		}

		missions, err := s.deps.Storage.UpsertMissions(ctx, []*models.Mission{mission})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to upsert mission: %v", err)
		}
		mission = missions[0]

		event := map[string]interface{}{
			"event_id":   fmt.Sprintf("mission-created-%d", mission.ID),
			"mission_id": mission.ID,
			"base_id":    mission.LaunchBaseID,
			"payload":    missionPb,
			"timestamp":  time.Now().Format(time.RFC3339),
		}
		eventBytes, _ := json.Marshal(event)
		err = s.deps.KafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(fmt.Sprintf("%d", mission.ID)),
			Value: eventBytes,
		})
		if err != nil {
			slog.Error("не удалось отправить missions.created", "mission_id", mission.ID, "err", err)
		}
	}

	return &mission_api.UpsertMissionsResponse{}, nil
}

func (s *MissionServiceImpl) GetMission(ctx context.Context, req *mission_api.GetMissionRequest) (*mission_api.GetMissionResponse, error) {
	missions, err := s.deps.Storage.GetMissionsByIDs(ctx, []uint64{req.MissionId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get mission: %v", err)
	}
	if len(missions) == 0 {
		return nil, status.Errorf(codes.NotFound, "mission not found")
	}

	mission := missions[0]
	pbMission := &models_pb.Mission{
		Id:        mission.ID,
		OpId:      mission.OperatorID,
		BaseId:    mission.LaunchBaseID,
		Status:    mission.Status,
		Lat:       mission.DestinationLat,
		Lon:       mission.DestinationLon,
		Alt:       mission.DestinationAlt,
		Payload:   mission.PayloadKg,
		CreatedAt: mission.CreatedAt.Format(time.RFC3339),
		Operator:  &models_pb.Operator{Id: mission.Operator.ID, Email: mission.Operator.Email, Name: mission.Operator.Name, CreatedAt: mission.Operator.CreatedAt.Format(time.RFC3339)},
		Base:      &models_pb.LaunchBase{Id: mission.LaunchBase.ID, Name: mission.LaunchBase.Name, Lat: mission.LaunchBase.Lat, Lon: mission.LaunchBase.Lon, Alt: mission.LaunchBase.Alt, CreatedAt: mission.LaunchBase.CreatedAt.Format(time.RFC3339)},
		Drones:    []*models_pb.Drone{},
	}

	return &mission_api.GetMissionResponse{Mission: pbMission}, nil
}
