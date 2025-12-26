package telemetry_processor

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (p *TelemetryProcessorImpl) Handle(ctx context.Context, telemetry *models.DroneTelemetry) error {
	return p.telemetrySvc.SaveTelemetry(ctx, telemetry)
}
