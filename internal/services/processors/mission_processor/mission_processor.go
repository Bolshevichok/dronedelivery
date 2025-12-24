package mission_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/services/droneService"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
)

type MissionProcessor interface {
	ProcessMissionCreated(ctx context.Context, missionID uint64)
}

type MissionProcessorImpl struct {
	missionService missionService.MissionService
	droneService   droneService.DroneService
}

func NewMissionProcessor(missionService missionService.MissionService, droneService droneService.DroneService) *MissionProcessorImpl {
	return &MissionProcessorImpl{
		missionService: missionService,
		droneService:   droneService,
	}
}

func (p *MissionProcessorImpl) ProcessMissionCreated(ctx context.Context, missionID uint64) {
	p.droneService.ProcessMissionCreated(ctx, missionID)
}
