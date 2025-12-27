package mission_api

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
)

type missionSvc interface {
	UpsertMissions(ctx context.Context, missions []*models.Mission) ([]*models.Mission, error)
	GetMission(ctx context.Context, missionID uint64) (*models.Mission, error)
	CreateOperator(ctx context.Context, operator *models.Operator) (uint64, error)
	CreateLaunchBase(ctx context.Context, launchBase *models.LaunchBase) (uint64, error)
	CreateDrone(ctx context.Context, drone *models.Drone) (uint64, error)
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
