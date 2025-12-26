package mission_lifecycle_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (p *MissionLifecycleProcessorImpl) Handle(ctx context.Context, mission *models.Mission) error {
	return p.missionSvc.ProcessMissionLifecycle(ctx, mission)
}
