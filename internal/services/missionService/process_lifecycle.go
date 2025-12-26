package missionService

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (s *MissionService) ProcessMissionLifecycle(ctx context.Context, event *models.MissionLifecycleEvent) error {
	if event == nil {
		return nil
	}

	return s.missionStorage.UpdateMissionStatus(ctx, event.MissionID, event.Status)
}
