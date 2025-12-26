package mission_api

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
	models_pb "github.com/Bolshevichok/dronedelivery/internal/pb/models"
)

func (a *MissionAPI) GetMission(ctx context.Context, req *mission_api.GetMissionRequest) (*mission_api.GetMissionResponse, error) {
	mission, err := a.missionSvc.GetMission(ctx, req.MissionId)
	if err != nil {
		return nil, err
	}
	return &mission_api.GetMissionResponse{Mission: mapMissionToResponse(mission)}, nil
}

func mapMissionToResponse(mission *models.Mission) *models_pb.Mission {
	return &models_pb.Mission{
		Id:        mission.ID,
		OpId:      mission.OperatorID,
		BaseId:    mission.LaunchBaseID,
		Status:    mission.Status,
		Lat:       mission.DestinationLat,
		Lon:       mission.DestinationLon,
		Alt:       mission.DestinationAlt,
		Payload:   mission.PayloadKg,
		CreatedAt: mission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Operator: &models_pb.Operator{
			Id:        mission.Operator.ID,
			Email:     mission.Operator.Email,
			Name:      mission.Operator.Name,
			CreatedAt: mission.Operator.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Base: &models_pb.LaunchBase{
			Id:        mission.LaunchBase.ID,
			Name:      mission.LaunchBase.Name,
			Lat:       mission.LaunchBase.Lat,
			Lon:       mission.LaunchBase.Lon,
			Alt:       mission.LaunchBase.Alt,
			CreatedAt: mission.LaunchBase.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
		Drones: []*models_pb.Drone{},
	}
}
