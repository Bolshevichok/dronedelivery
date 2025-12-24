package telemetry_consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	telemetryprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/telemetry_processor"
	"github.com/segmentio/kafka-go"
)

type TelemetryConsumer interface {
	Start(ctx context.Context) error
	Stop() error
}

type TelemetryConsumerImpl struct {
	processor *telemetryprocessor.TelemetryProcessorImpl
	reader    *kafka.Reader
	cancel    context.CancelFunc
}

func NewTelemetryConsumer(processor *telemetryprocessor.TelemetryProcessorImpl, brokers []string, topic string) *TelemetryConsumerImpl {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "telemetry-extractor",
	})

	return &TelemetryConsumerImpl{
		processor: processor,
		reader:    reader,
	}
}

func (c *TelemetryConsumerImpl) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	slog.Info("Starting telemetry consumer")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				slog.Error("Failed to read message", "error", err)
				continue
			}

			var telemetry map[string]interface{}
			if err := json.Unmarshal(msg.Value, &telemetry); err != nil {
				slog.Error("Failed to unmarshal telemetry", "error", err)
				continue
			}

			slog.Info("Processing telemetry", "drone_id", telemetry["drone_id"])

			if err := c.processor.ProcessTelemetry(ctx, telemetry); err != nil {
				slog.Error("Failed to process telemetry", "error", err)
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				slog.Error("Failed to commit message", "error", err)
			}
		}
	}
}

func (c *TelemetryConsumerImpl) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}
