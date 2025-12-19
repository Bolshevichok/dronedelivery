package missionService

import (
	"context"
	"time"

	missionv1 "github.com/Bolshevichok/dronedelivery/internal/pb/mission/v1"
	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Dependencies for the service
type Dependencies struct {
	Storage *pgstorage.PGstorage
	// Add Kafka producer, etc.
}

// MissionService implements the gRPC MissionServiceServer
type MissionService struct {
	missionv1.UnimplementedMissionServiceServer
	deps *Dependencies
}

func NewMissionService(deps *Dependencies) *MissionService {
	return &MissionService{deps: deps}
}

// CreateMission creates a new mission
func (s *MissionService) CreateMission(ctx context.Context, req *missionv1.CreateMissionRequest) (*missionv1.CreateMissionResponse, error) {
	mission := &pgstorage.Mission{
		OperatorID:     req.OperatorId,
		LaunchBaseID:   req.LaunchBaseId,
		Status:         req.Status,
		DestinationLat: req.DestinationLat,
		DestinationLon: req.DestinationLon,
		DestinationAlt: req.DestinationAlt,
		PayloadKg:      req.PayloadKg,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}

	// TODO: Implement UpsertMissions in pgstorage
	// For now, assume ID is assigned
	mission.ID = 1 // Placeholder

	return &missionv1.CreateMissionResponse{MissionId: mission.ID}, nil
}

// GetMission retrieves a mission by ID
func (s *MissionService) GetMission(ctx context.Context, req *missionv1.GetMissionRequest) (*missionv1.GetMissionResponse, error) {
	missions, err := s.deps.Storage.GetMissionsByIDs(ctx, []uint64{req.MissionId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get mission: %v", err)
	}
	if len(missions) == 0 {
		return nil, status.Errorf(codes.NotFound, "mission not found")
	}

	mission := missions[0]
	// Convert to protobuf
	pbMission := &missionv1.Mission{
		Id:             mission.ID,
		OperatorId:     mission.OperatorID,
		LaunchBaseId:   mission.LaunchBaseID,
		Status:         mission.Status,
		DestinationLat: mission.DestinationLat,
		DestinationLon: mission.DestinationLon,
		DestinationAlt: mission.DestinationAlt,
		PayloadKg:      mission.PayloadKg,
		CreatedAt:      mission.CreatedAt,
		// TODO: Load related Operator, LaunchBase, Drones
	}

	return &missionv1.GetMissionResponse{Mission: pbMission}, nil
}

// ListMissions lists all missions (basic implementation)
func (s *MissionService) ListMissions(ctx context.Context, req *missionv1.ListMissionsRequest) (*missionv1.ListMissionsResponse, error) {
	// TODO: Implement proper listing with filters
	// For now, return empty
	return &missionv1.ListMissionsResponse{Missions: []*missionv1.Mission{}}, nil
}

// GetMissionTelemetry placeholder
func (s *MissionService) GetMissionTelemetry(ctx context.Context, req *missionv1.GetMissionTelemetryRequest) (*missionv1.GetMissionTelemetryResponse, error) {
	// TODO: Implement telemetry retrieval from Redis
	return &missionv1.GetMissionTelemetryResponse{Telemetry: &missionv1.Telemetry{}}, nil
}

// WatchMission placeholder for streaming
func (s *MissionService) WatchMission(req *missionv1.WatchMissionRequest, stream missionv1.MissionService_WatchMissionServer) error {
	// TODO: Implement streaming updates
	return nil
}
