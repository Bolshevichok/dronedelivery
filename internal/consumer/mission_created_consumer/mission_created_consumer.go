package mission_created_consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	missionprocessor "github.com/Bolshevichok/dronedelivery/internal/services/processors/mission_processor"
	"github.com/segmentio/kafka-go"
)

type MissionCreatedConsumer interface {
	Start(ctx context.Context) error
	Stop() error
}

type MissionCreatedConsumerImpl struct {
	processor *missionprocessor.MissionProcessorImpl
	reader    *kafka.Reader
	cancel    context.CancelFunc
}

func NewMissionCreatedConsumer(processor *missionprocessor.MissionProcessorImpl, brokers []string, topic string) *MissionCreatedConsumerImpl {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "mission-created-consumer",
	})

	return &MissionCreatedConsumerImpl{
		processor: processor,
		reader:    reader,
	}
}

func (c *MissionCreatedConsumerImpl) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	slog.Info("Starting mission created consumer")

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

			var event map[string]interface{}
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				slog.Error("Failed to unmarshal event", "error", err)
				continue
			}

			missionID, ok := event["mission_id"].(float64)
			if !ok {
				slog.Error("Invalid mission_id in event")
				continue
			}

			slog.Info("Processing mission created event", "mission_id", uint64(missionID))

			c.processor.ProcessMissionCreated(ctx, uint64(missionID))

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				slog.Error("Failed to commit message", "error", err)
			}
		}
	}
}

func (c *MissionCreatedConsumerImpl) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}
