package mission_lifecycle_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

type missionSvc interface {
	ProcessMissionLifecycle(ctx context.Context, event *models.MissionLifecycleEvent) error
}

type MissionLifecycleProcessor interface {
	ProcessMissionLifecycle(ctx context.Context, event *models.MissionLifecycleEvent) error
	Handle(ctx context.Context, event *models.MissionLifecycleEvent) error
}

type MissionLifecycleProcessorImpl struct {
	missionSvc missionSvc
}

func NewMissionLifecycleProcessor(missionSvc missionSvc) *MissionLifecycleProcessorImpl {
	return &MissionLifecycleProcessorImpl{missionSvc: missionSvc}
}

func (p *MissionLifecycleProcessorImpl) ProcessMissionLifecycle(ctx context.Context, event *models.MissionLifecycleEvent) error {
	return p.missionSvc.ProcessMissionLifecycle(ctx, event)
}
