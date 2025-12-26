package telemetry_processor

import (
	"context"
)

func (p *TelemetryProcessorImpl) Handle(ctx context.Context, telemetry map[string]interface{}) error {
	return p.telemetrySvc.SaveTelemetry(ctx, telemetry)
}
