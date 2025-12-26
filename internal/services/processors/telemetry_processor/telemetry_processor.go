package telemetry_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

type telemetrySvc interface {
	SaveTelemetry(ctx context.Context, telemetry *models.DroneTelemetry) error
}

type TelemetryProcessor interface {
	ProcessTelemetry(ctx context.Context, telemetry *models.DroneTelemetry) error
	Handle(ctx context.Context, telemetry *models.DroneTelemetry) error
}

type TelemetryProcessorImpl struct {
	telemetrySvc telemetrySvc
}

func NewTelemetryProcessor(telemetrySvc telemetrySvc) *TelemetryProcessorImpl {
	return &TelemetryProcessorImpl{
		telemetrySvc: telemetrySvc,
	}
}

func (p *TelemetryProcessorImpl) ProcessTelemetry(ctx context.Context, telemetry *models.DroneTelemetry) error {
	return p.telemetrySvc.SaveTelemetry(ctx, telemetry)
}
