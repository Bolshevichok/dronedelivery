package mission_processor

import (
	"context"
)

type droneSvc interface {
	ProcessMissionCreated(ctx context.Context, missionID uint64) error
}

type MissionProcessor interface {
	ProcessMissionCreated(ctx context.Context, missionID uint64) error
	Handle(ctx context.Context, missionID uint64) error
}

type MissionProcessorImpl struct {
	droneSvc droneSvc
}

func NewMissionProcessor(droneSvc droneSvc) *MissionProcessorImpl {
	return &MissionProcessorImpl{
		droneSvc: droneSvc,
	}
}

func (p *MissionProcessorImpl) ProcessMissionCreated(ctx context.Context, missionID uint64) error {
	return p.droneSvc.ProcessMissionCreated(ctx, missionID)
}
