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

	mission := missions[0]

	// Read telemetry for drones in the mission
	for _, drone := range mission.Drones {
		telemetry, err := s.GetDroneTelemetry(ctx, fmt.Sprintf("%d", drone.ID))
		if err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to get telemetry for drone %d: %v\n", drone.ID, err)
			continue
		}
		// For example, update drone status or add telemetry to model
		// Here we just print for demo
		fmt.Printf("Telemetry for drone %d: %s\n", drone.ID, telemetry)
		// If needed, update mission.Drones[i].Status based on telemetry
	}

	return mission, nil
}
