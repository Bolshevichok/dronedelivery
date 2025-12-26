package telemetryService

import (
	"context"
	"fmt"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/go-redis/redis/v8"
)

// TelemetryService interface
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
