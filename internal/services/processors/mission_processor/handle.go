package mission_processor

import (
	"context"
)

func (p *MissionProcessorImpl) Handle(ctx context.Context, missionID uint64) error {
	return p.droneSvc.ProcessMissionCreated(ctx, missionID)
}
