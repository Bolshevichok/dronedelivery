package mission_api

import (
	"context"
	"time"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
)

func (a *MissionAPI) CreateDrone(ctx context.Context, req *mission_api.CreateDroneRequest) (*mission_api.CreateDroneResponse, error) {
	id, err := a.missionSvc.CreateDrone(ctx, mapCreateDroneRequest(req))
	if err != nil {
		return nil, err
	}
	return &mission_api.CreateDroneResponse{Id: id}, nil
}

func mapCreateDroneRequest(req *mission_api.CreateDroneRequest) *models.Drone {
	return &models.Drone{
		Serial:       req.Drone.Serial,
		Model:        req.Drone.Model,
		Status:       req.Drone.Status,
		LaunchBaseID: req.Drone.BaseId,
		CreatedAt:    time.Now(),
	}
}
