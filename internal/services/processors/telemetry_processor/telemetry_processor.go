package telemetry_processor

import (
	"context"
)

type telemetrySvc interface {
	SaveTelemetry(ctx context.Context, telemetry map[string]interface{}) error
}

type TelemetryProcessor interface {
	ProcessTelemetry(ctx context.Context, telemetry map[string]interface{}) error
	Handle(ctx context.Context, telemetry map[string]interface{}) error
}

type TelemetryProcessorImpl struct {
	telemetrySvc telemetrySvc
}

func NewTelemetryProcessor(telemetrySvc telemetrySvc) *TelemetryProcessorImpl {
	return &TelemetryProcessorImpl{
		telemetrySvc: telemetrySvc,
	}
}

func (p *TelemetryProcessorImpl) ProcessTelemetry(ctx context.Context, telemetry map[string]interface{}) error {
	return p.telemetrySvc.SaveTelemetry(ctx, telemetry)
}
