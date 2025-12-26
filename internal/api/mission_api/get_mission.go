package mission_api

import (
	"context"
	"time"

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
	if mission == nil {
		return nil
	}

	return &models_pb.Mission{
		Id:        mission.ID,
		OpId:      mission.OperatorID,
		BaseId:    mission.LaunchBaseID,
		Status:    mission.Status,
		Lat:       mission.DestinationLat,
		Lon:       mission.DestinationLon,
		Alt:       mission.DestinationAlt,
		Payload:   mission.PayloadKg,
		CreatedAt: formatTimeRFC3339(mission.CreatedAt),
	}
}

func formatTimeRFC3339(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
