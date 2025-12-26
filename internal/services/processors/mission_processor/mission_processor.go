package mission_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

type droneSvc interface {
	ProcessMissionCreated(ctx context.Context, mission *models.Mission) error
}

type MissionProcessor interface {
	ProcessMissionCreated(ctx context.Context, mission *models.Mission) error
	Handle(ctx context.Context, mission *models.Mission) error
}

type MissionProcessorImpl struct {
	droneSvc droneSvc
}

func NewMissionProcessor(droneSvc droneSvc) *MissionProcessorImpl {
	return &MissionProcessorImpl{
		droneSvc: droneSvc,
	}
}

func (p *MissionProcessorImpl) ProcessMissionCreated(ctx context.Context, mission *models.Mission) error {
	return p.droneSvc.ProcessMissionCreated(ctx, mission)
}
