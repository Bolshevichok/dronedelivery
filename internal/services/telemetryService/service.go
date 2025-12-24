package telemetryService

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/go-redis/redis/v8"
)

type TelemetryService interface {
	SaveTelemetry(ctx context.Context, telemetry map[string]interface{}) error
}

type TelemetryServiceImpl struct {
	redisClient *redis.Client
}

func NewTelemetryService(cfg *config.Config) TelemetryService {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})

	return &TelemetryServiceImpl{
		redisClient: redisClient,
	}
}

func (s *TelemetryServiceImpl) SaveTelemetry(ctx context.Context, telemetry map[string]interface{}) error {
	droneID, ok := telemetry["drone_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid drone_id in telemetry")
	}

	key := fmt.Sprintf("telemetry:%d", uint64(droneID))

	data, err := json.Marshal(telemetry)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, key, string(data), time.Hour).Err()
}
