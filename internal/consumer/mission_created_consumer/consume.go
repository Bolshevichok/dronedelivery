package mission_created_consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

func (c *MissionCreatedConsumerImpl) Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           "MissionService_group",
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error("MissionCreatedConsumer.Consume error", "error", err.Error())
		}
		var event map[string]interface{}
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			slog.Error("parse", "error", err)
			continue
		}
		missionID, ok := event["mission_id"].(float64)
		if !ok {
			slog.Error("Invalid mission_id in event")
			continue
		}
		err = c.processor.ProcessMissionCreated(ctx, uint64(missionID))
		if err != nil {
			slog.Error("ProcessMissionCreated error", "error", err.Error())
		}
	}
}
