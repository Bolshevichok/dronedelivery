package mission_lifecycle_consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	"github.com/segmentio/kafka-go"
)

func (c *MissionLifecycleConsumerImpl) Consume(ctx context.Context) {
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
			slog.Error("MissionLifecycleConsumer.Consume error", "error", err.Error())
			continue
		}

		var event *models.MissionLifecycleEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			slog.Error("parse", "error", err)
			continue
		}
		if event == nil || event.MissionID == 0 || event.DroneID == 0 || event.Status == "" {
			slog.Error("Invalid mission lifecycle")
			continue
		}

		err = c.processor.ProcessMissionLifecycle(ctx, event)
		if err != nil {
			slog.Error("ProcessMissionLifecycle error", "error", err.Error())
		}
	}
}
