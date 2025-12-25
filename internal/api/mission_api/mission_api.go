package mission_api

import (
	"context"

	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
)

type missionSvc interface {
	UpsertMissions(ctx context.Context, req *mission_api.UpsertMissionsRequest) (*mission_api.UpsertMissionsResponse, error)
	GetMission(ctx context.Context, req *mission_api.GetMissionRequest) (*mission_api.GetMissionResponse, error)
}

type MissionAPI struct {
	mission_api.UnimplementedMissionServiceServer
	missionSvc missionSvc
}

func NewMissionAPI(missionSvc missionSvc) *MissionAPI {
	return &MissionAPI{
		missionSvc: missionSvc,
	}
}

func (a *MissionAPI) UpsertMissions(ctx context.Context, req *mission_api.UpsertMissionsRequest) (*mission_api.UpsertMissionsResponse, error) {
	return a.missionSvc.UpsertMissions(ctx, req)
}

func (a *MissionAPI) GetMission(ctx context.Context, req *mission_api.GetMissionRequest) (*mission_api.GetMissionResponse, error) {
	return a.missionSvc.GetMission(ctx, req)
}
