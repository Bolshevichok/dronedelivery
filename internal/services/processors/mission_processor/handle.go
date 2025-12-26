package mission_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (p *MissionProcessorImpl) Handle(ctx context.Context, mission *models.MissionInfo) error {
	return p.droneSvc.ProcessMissionCreated(ctx, mission)
}
