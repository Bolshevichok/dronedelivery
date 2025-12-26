package bootstrap

import (
	"github.com/Bolshevichok/dronedelivery/internal/services/droneService"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
	missionlifecycleprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/mission_lifecycle_processor"
	missionprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/mission_processor"
	telemetryprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/telemetry_processor"
	"github.com/Bolshevichok/dronedelivery/internal/services/telemetryService"
)

func InitMissionProcessor(droneService *droneService.DroneService) *missionprocessor.MissionProcessorImpl {
	return missionprocessor.NewMissionProcessor(droneService)
}

func InitTelemetryProcessor(telemetryService telemetryService.TelemetryService) *telemetryprocessor.TelemetryProcessorImpl {
	return telemetryprocessor.NewTelemetryProcessor(telemetryService)
}

func InitMissionLifecycleProcessor(missionService *missionService.MissionService) *missionlifecycleprocessor.MissionLifecycleProcessorImpl {
	return missionlifecycleprocessor.NewMissionLifecycleProcessor(missionService)
}
