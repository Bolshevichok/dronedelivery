package mission_api

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
	models_pb "github.com/Bolshevichok/dronedelivery/internal/pb/models"
	"github.com/samber/lo"
)

func (a *MissionAPI) UpsertMissions(ctx context.Context, req *mission_api.UpsertMissionsRequest) (*mission_api.UpsertMissionsResponse, error) {
	_, err := a.missionSvc.UpsertMissions(ctx, mapUpsertMissionsRequest(req))
	if err != nil {
		return nil, err
	}
	return &mission_api.UpsertMissionsResponse{}, nil
}

func mapUpsertMissionsRequest(req *mission_api.UpsertMissionsRequest) []*models.Mission {
	return lo.Map(req.Missions, func(m *models_pb.Mission, _ int) *models.Mission {
		return &models.Mission{
			OperatorID:     m.OpId,
			LaunchBaseID:   m.BaseId,
			Status:         m.Status,
			DestinationLat: m.Lat,
			DestinationLon: m.Lon,
			DestinationAlt: m.Alt,
			PayloadKg:      m.Payload,
		}
	})
}
