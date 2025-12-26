package bootstrap

import (
	"github.com/Bolshevichok/dronedelivery/internal/api/mission_api"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
)

func InitMissionAPI(missionService *missionService.MissionService) *mission_api.MissionAPI {
	return mission_api.NewMissionAPI(missionService)
}
