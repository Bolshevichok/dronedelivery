package telemetryService

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

func (s *TelemetryServiceImpl) SaveTelemetry(ctx context.Context, telemetry *models.DroneTelemetry) error {
	if telemetry == nil || telemetry.DroneID == 0 {
		return fmt.Errorf("invalid drone_id in telemetry")
	}

	key := fmt.Sprintf("telemetry:%d", telemetry.DroneID)

	data, err := json.Marshal(telemetry)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, key, string(data), time.Hour).Err()
}
