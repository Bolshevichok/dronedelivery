package telemetry_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/services/telemetryService"
)

type TelemetryProcessor interface {
	ProcessTelemetry(ctx context.Context, telemetry map[string]interface{}) error
}

type TelemetryProcessorImpl struct {
	telemetryService telemetryService.TelemetryService
}

func NewTelemetryProcessor(telemetryService telemetryService.TelemetryService) *TelemetryProcessorImpl {
	return &TelemetryProcessorImpl{
		telemetryService: telemetryService,
	}
}

func (p *TelemetryProcessorImpl) ProcessTelemetry(ctx context.Context, telemetry map[string]interface{}) error {
	return p.telemetryService.SaveTelemetry(ctx, telemetry)
}
