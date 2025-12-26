package telemetry_consumer

import (
	"context"
)

type telemetryProcessor interface {
	ProcessTelemetry(ctx context.Context, telemetry map[string]interface{}) error
}

type TelemetryConsumer interface {
	Consume(ctx context.Context)
}

type TelemetryConsumerImpl struct {
	processor   telemetryProcessor
	kafkaBroker []string
	topicName   string
}

func NewTelemetryConsumer(processor telemetryProcessor, kafkaBroker []string, topicName string) *TelemetryConsumerImpl {
	return &TelemetryConsumerImpl{
		processor:   processor,
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
	}
}
