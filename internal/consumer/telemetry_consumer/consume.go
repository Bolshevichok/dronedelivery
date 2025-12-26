package telemetry_consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

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
		}
		var telemetry map[string]interface{}
		err = json.Unmarshal(msg.Value, &telemetry)
		if err != nil {
			slog.Error("parse", "error", err)
			continue
		}
		err = c.processor.ProcessTelemetry(ctx, telemetry)
		if err != nil {
			slog.Error("ProcessTelemetry", "error", err)
		}
	}
}
