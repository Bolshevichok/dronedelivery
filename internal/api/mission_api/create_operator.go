package mission_api

import (
	"context"
	"time"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
)

func (a *MissionAPI) CreateOperator(ctx context.Context, req *mission_api.CreateOperatorRequest) (*mission_api.CreateOperatorResponse, error) {
	id, err := a.missionSvc.CreateOperator(ctx, mapCreateOperatorRequest(req))
	if err != nil {
		return nil, err
	}
	return &mission_api.CreateOperatorResponse{Id: id}, nil
}

func mapCreateOperatorRequest(req *mission_api.CreateOperatorRequest) *models.Operator {
	return &models.Operator{
		Email:     req.Operator.Email,
		Name:      req.Operator.Name,
		CreatedAt: time.Now(),
	}
}
