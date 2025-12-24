package bootstrap

import (
	"github.com/Bolshevichok/dronedelivery/internal/services/droneService"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
	missionprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/mission_processor"
	telemetryprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/telemetry_processor"
	"github.com/Bolshevichok/dronedelivery/internal/services/telemetryService"
)

func InitMissionProcessor(missionService missionService.MissionService, droneService droneService.DroneService) *missionprocessor.MissionProcessorImpl {
	return missionprocessor.NewMissionProcessor(missionService, droneService)
}

func InitTelemetryProcessor(telemetryService telemetryService.TelemetryService) *telemetryprocessor.TelemetryProcessorImpl {
	return telemetryprocessor.NewTelemetryProcessor(telemetryService)
}
