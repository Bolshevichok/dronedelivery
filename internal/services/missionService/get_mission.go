package missionService

import (
	"context"
	"time"

	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
	models_pb "github.com/Bolshevichok/dronedelivery/internal/pb/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MissionService) GetMission(ctx context.Context, req *mission_api.GetMissionRequest) (*mission_api.GetMissionResponse, error) {
	missions, err := s.missionStorage.GetMissionsByIDs(ctx, []uint64{req.MissionId})
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
