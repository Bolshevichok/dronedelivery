package mission_api

import (
	"context"
	"time"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
)

func (a *MissionAPI) CreateLaunchBase(ctx context.Context, req *mission_api.CreateLaunchBaseRequest) (*mission_api.CreateLaunchBaseResponse, error) {
	id, err := a.missionSvc.CreateLaunchBase(ctx, mapCreateLaunchBaseRequest(req))
	if err != nil {
		return nil, err
	}
	return &mission_api.CreateLaunchBaseResponse{Id: id}, nil
}

func mapCreateLaunchBaseRequest(req *mission_api.CreateLaunchBaseRequest) *models.LaunchBase {
	return &models.LaunchBase{
		Name:      req.LaunchBase.Name,
		Lat:       req.LaunchBase.Lat,
		Lon:       req.LaunchBase.Lon,
		Alt:       req.LaunchBase.Alt,
		CreatedAt: time.Now(),
	}
}
