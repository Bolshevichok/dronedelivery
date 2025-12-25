package mission_api

import (
	"context"

	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
)

func (a *MissionAPI) GetMission(ctx context.Context, req *mission_api.GetMissionRequest) (*mission_api.GetMissionResponse, error) {
	return a.missionSvc.GetMission(ctx, req)
}
