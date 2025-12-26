package mission_lifecycle_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (p *MissionLifecycleProcessorImpl) Handle(ctx context.Context, event *models.MissionLifecycleEvent) error {
	return p.missionSvc.ProcessMissionLifecycle(ctx, event)
}
