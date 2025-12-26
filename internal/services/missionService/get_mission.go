package missionService

import (
	"context"
	"fmt"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (s *MissionService) GetMission(ctx context.Context, missionID uint64) (*models.Mission, error) {
	missions, err := s.missionStorage.GetMissionsByIDs(ctx, []uint64{missionID})
	if err != nil {
		return nil, err
	}
	if len(missions) == 0 {
		return nil, fmt.Errorf("mission not found")
	}

	return missions[0], nil
}
