package mission_api

import (
	"context"

	missionv1 "github.com/Bolshevichok/dronedelivery/internal/pb/mission/v1"
)

type missionSvc interface {
	CreateMission(ctx context.Context, req *missionv1.CreateMissionRequest) (*missionv1.CreateMissionResponse, error)
	GetMission(ctx context.Context, req *missionv1.GetMissionRequest) (*missionv1.GetMissionResponse, error)
	ListMissions(ctx context.Context, req *missionv1.ListMissionsRequest) (*missionv1.ListMissionsResponse, error)
	GetMissionTelemetry(ctx context.Context, req *missionv1.GetMissionTelemetryRequest) (*missionv1.GetMissionTelemetryResponse, error)
	WatchMission(req *missionv1.WatchMissionRequest, stream missionv1.MissionService_WatchMissionServer) error
}

type MissionAPI struct {
	missionv1.UnimplementedMissionServiceServer
	missionSvc missionSvc
}

func NewMissionAPI(missionSvc missionSvc) *MissionAPI {
	return &MissionAPI{
		missionSvc: missionSvc,
	}
}

func (a *MissionAPI) CreateMission(ctx context.Context, req *missionv1.CreateMissionRequest) (*missionv1.CreateMissionResponse, error) {
	return a.missionSvc.CreateMission(ctx, req)
}

func (a *MissionAPI) GetMission(ctx context.Context, req *missionv1.GetMissionRequest) (*missionv1.GetMissionResponse, error) {
	return a.missionSvc.GetMission(ctx, req)
}

func (a *MissionAPI) ListMissions(ctx context.Context, req *missionv1.ListMissionsRequest) (*missionv1.ListMissionsResponse, error) {
	return a.missionSvc.ListMissions(ctx, req)
}

func (a *MissionAPI) GetMissionTelemetry(ctx context.Context, req *missionv1.GetMissionTelemetryRequest) (*missionv1.GetMissionTelemetryResponse, error) {
	return a.missionSvc.GetMissionTelemetry(ctx, req)
}

func (a *MissionAPI) WatchMission(req *missionv1.WatchMissionRequest, stream missionv1.MissionService_WatchMissionServer) error {
	return a.missionSvc.WatchMission(req, stream)
}
