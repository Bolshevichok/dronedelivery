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
	missionv1 "github.com/Bolshevichok/dronedelivery/internal/pb/mission/v1"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MissionService interface
type MissionService interface {
	CreateMission(ctx context.Context, req *missionv1.CreateMissionRequest) (*missionv1.CreateMissionResponse, error)
	GetMission(ctx context.Context, req *missionv1.GetMissionRequest) (*missionv1.GetMissionResponse, error)
	ListMissions(ctx context.Context, req *missionv1.ListMissionsRequest) (*missionv1.ListMissionsResponse, error)
	GetMissionTelemetry(ctx context.Context, req *missionv1.GetMissionTelemetryRequest) (*missionv1.GetMissionTelemetryResponse, error)
	WatchMission(req *missionv1.WatchMissionRequest, stream missionv1.MissionService_WatchMissionServer) error
}

type Dependencies struct {
	Storage     internal.Storage
	KafkaWriter *kafka.Writer
	RedisClient *redis.Client
}

type MissionServiceImpl struct {
	missionv1.UnimplementedMissionServiceServer
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

func (s *MissionServiceImpl) CreateMission(ctx context.Context, req *missionv1.CreateMissionRequest) (*missionv1.CreateMissionResponse, error) {
	mission := &models.Mission{
		OperatorID:     req.OperatorId,
		LaunchBaseID:   req.LaunchBaseId,
		Status:         "created",
		DestinationLat: req.DestinationLat,
		DestinationLon: req.DestinationLon,
		DestinationAlt: req.DestinationAlt,
		PayloadKg:      req.PayloadKg,
		CreatedAt:      time.Now(),
	}

	missions, err := s.deps.Storage.UpsertMissions(ctx, []*models.Mission{mission})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create mission: %v", err)
	}
	mission = missions[0]

	event := map[string]interface{}{
		"event_id":   fmt.Sprintf("mission-created-%d", mission.ID),
		"mission_id": mission.ID,
		"base_id":    mission.LaunchBaseID,
		"payload":    req,
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

	return &missionv1.CreateMissionResponse{MissionId: mission.ID}, nil
}

func (s *MissionServiceImpl) GetMission(ctx context.Context, req *missionv1.GetMissionRequest) (*missionv1.GetMissionResponse, error) {
	missions, err := s.deps.Storage.GetMissionsByIDs(ctx, []uint64{req.MissionId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get mission: %v", err)
	}
	if len(missions) == 0 {
		return nil, status.Errorf(codes.NotFound, "mission not found")
	}

	mission := missions[0]
	pbMission := &missionv1.Mission{
		Id:             mission.ID,
		OperatorId:     mission.OperatorID,
		LaunchBaseId:   mission.LaunchBaseID,
		Status:         mission.Status,
		DestinationLat: mission.DestinationLat,
		DestinationLon: mission.DestinationLon,
		DestinationAlt: mission.DestinationAlt,
		PayloadKg:      mission.PayloadKg,
		CreatedAt:      mission.CreatedAt.Format(time.RFC3339),
	}

	return &missionv1.GetMissionResponse{Mission: pbMission}, nil
}

func (s *MissionServiceImpl) ListMissions(ctx context.Context, req *missionv1.ListMissionsRequest) (*missionv1.ListMissionsResponse, error) {
	return &missionv1.ListMissionsResponse{Missions: []*missionv1.Mission{}}, nil
}

func (s *MissionServiceImpl) GetMissionTelemetry(ctx context.Context, req *missionv1.GetMissionTelemetryRequest) (*missionv1.GetMissionTelemetryResponse, error) {
	missionDrones, err := s.deps.Storage.GetMissionDronesByMissionIDs(ctx, []uint64{req.MissionId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get mission drones: %v", err)
	}
	if len(missionDrones) == 0 {
		return &missionv1.GetMissionTelemetryResponse{Telemetry: &missionv1.Telemetry{}}, nil
	}
	droneID := missionDrones[0].DroneID

	key := fmt.Sprintf("telemetry:%d", droneID)
	val, err := s.deps.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return &missionv1.GetMissionTelemetryResponse{Telemetry: &missionv1.Telemetry{}}, nil
	}

	type telemetryDTO struct {
		DroneID   uint64  `json:"drone_id"`
		MissionID uint64  `json:"mission_id"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		Alt       float64 `json:"alt"`
		Timestamp string  `json:"timestamp"`
	}

	var telemetry telemetryDTO
	if err := json.Unmarshal([]byte(val), &telemetry); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal telemetry: %v", err)
	}

	pbTelemetry := &missionv1.Telemetry{
		DroneId:   telemetry.DroneID,
		Lat:       telemetry.Lat,
		Lon:       telemetry.Lon,
		Alt:       telemetry.Alt,
		Timestamp: telemetry.Timestamp,
	}

	return &missionv1.GetMissionTelemetryResponse{Telemetry: pbTelemetry}, nil
}

func (s *MissionServiceImpl) WatchMission(req *missionv1.WatchMissionRequest, stream missionv1.MissionService_WatchMissionServer) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			missions, err := s.deps.Storage.GetMissionsByIDs(stream.Context(), []uint64{req.MissionId})
			if err != nil || len(missions) == 0 {
				continue
			}
			mission := missions[0]

			telemetryResp, err := s.GetMissionTelemetry(stream.Context(), &missionv1.GetMissionTelemetryRequest{MissionId: req.MissionId})
			if err != nil {
				continue
			}

			update := &missionv1.MissionUpdate{
				Status:    mission.Status,
				Telemetry: telemetryResp.Telemetry,
				Timestamp: time.Now().Format(time.RFC3339),
			}

			if err := stream.Send(update); err != nil {
				return err
			}
		}
	}
}
