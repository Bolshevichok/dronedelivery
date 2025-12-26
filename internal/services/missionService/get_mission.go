package missionService

import (
	"context"
	"fmt"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (s *MissionService) GetMission(ctx context.Context, missionID uint64) (*models.Mission, error) {
	missions, err := s.missionStorage.GetMissionsByIDs(ctx, []uint64{missionID})
	if err != nil {
		return nil, fmt.Errorf("failed to get mission: %w", err)
	}
	if len(missions) == 0 {
		return nil, fmt.Errorf("mission not found")
	}
	mission := missions[0]

	// Load related data
	operators, err := s.missionStorage.GetOperatorsByIDs(ctx, []uint64{mission.OperatorID})
	if err != nil {
		return nil, fmt.Errorf("failed to get operator: %w", err)
	}
	if len(operators) > 0 {
		mission.Operator = *operators[0]
	}

	launchBases, err := s.missionStorage.GetLaunchBasesByIDs(ctx, []uint64{mission.LaunchBaseID})
	if err != nil {
		return nil, fmt.Errorf("failed to get launch base: %w", err)
	}
	if len(launchBases) > 0 {
		mission.LaunchBase = *launchBases[0]
	}

	missionDrones, err := s.missionStorage.GetMissionDronesByMissionID(ctx, missionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get mission drones: %w", err)
	}
	droneIDs := make([]uint64, len(missionDrones))
	for i, md := range missionDrones {
		droneIDs[i] = md.DroneID
	}
	drones, err := s.missionStorage.GetDronesByIDs(ctx, droneIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get drones: %w", err)
	}
	mission.Drones = make([]models.Drone, len(drones))
	for i, d := range drones {
		mission.Drones[i] = *d
	}

	return mission, nil
}
