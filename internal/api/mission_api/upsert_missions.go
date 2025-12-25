package mission_api

import (
	"context"

	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
)

func (a *MissionAPI) UpsertMissions(ctx context.Context, req *mission_api.UpsertMissionsRequest) (*mission_api.UpsertMissionsResponse, error) {
	return a.missionSvc.UpsertMissions(ctx, req)
}
