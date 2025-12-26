package mission_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (p *MissionProcessorImpl) Handle(ctx context.Context, mission *models.Mission) error {
	return p.droneSvc.ProcessMissionCreated(ctx, mission)
}
