package telemetry_consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	"github.com/segmentio/kafka-go"
)

func (c *TelemetryConsumerImpl) Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           "TelemetryService_group",
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error("TelemetryConsumer.Consume error", "error", err.Error())
			continue
		}
		var telemetry *models.DroneTelemetry
		err = json.Unmarshal(msg.Value, &telemetry)
		if err != nil {
			slog.Error("parse", "error", err)
			continue
		}
		if telemetry == nil || telemetry.DroneID == 0 {
			slog.Error("Invalid drone_id in telemetry")
			continue
		}
		err = c.processor.ProcessTelemetry(ctx, telemetry)
		if err != nil {
			slog.Error("ProcessTelemetry", "error", err)
		}
	}
}
