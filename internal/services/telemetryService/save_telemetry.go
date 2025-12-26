package telemetryService

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (s *TelemetryServiceImpl) SaveTelemetry(ctx context.Context, telemetry *models.DroneTelemetry) error {
	if telemetry == nil || telemetry.DroneID == 0 {
		return fmt.Errorf("invalid drone_id in telemetry")
	}

	data, err := json.Marshal(telemetry)
	if err != nil {
		return err
	}

	// Push to Redis list for consumption
	return s.redisClient.LPush(ctx, "telemetry_queue", string(data)).Err()
}
