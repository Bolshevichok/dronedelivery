package missionService

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (s *MissionService) ProcessMissionLifecycle(ctx context.Context, mission *models.Mission) error {
	if mission == nil {
		return nil
	}
	return s.missionStorage.UpdateMissionStatus(ctx, mission.ID, mission.Status)
}
