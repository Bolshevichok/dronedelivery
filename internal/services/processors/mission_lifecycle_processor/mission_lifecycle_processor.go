package mission_lifecycle_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

type missionSvc interface {
	ProcessMissionLifecycle(ctx context.Context, mission *models.Mission) error
}

type MissionLifecycleProcessor interface {
	ProcessMissionLifecycle(ctx context.Context, mission *models.Mission) error
	Handle(ctx context.Context, mission *models.Mission) error
}

type MissionLifecycleProcessorImpl struct {
	missionSvc missionSvc
}

func NewMissionLifecycleProcessor(missionSvc missionSvc) *MissionLifecycleProcessorImpl {
	return &MissionLifecycleProcessorImpl{missionSvc: missionSvc}
}

func (p *MissionLifecycleProcessorImpl) ProcessMissionLifecycle(ctx context.Context, mission *models.Mission) error {
	return p.missionSvc.ProcessMissionLifecycle(ctx, mission)
}
