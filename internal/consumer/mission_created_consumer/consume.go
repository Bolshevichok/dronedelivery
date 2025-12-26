package mission_created_consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	"github.com/segmentio/kafka-go"
)

func (c *MissionCreatedConsumerImpl) Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBroker,
		GroupID:           "DroneService_group",
		Topic:             c.topicName,
		HeartbeatInterval: 3 * time.Second,
		SessionTimeout:    30 * time.Second,
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error("MissionCreatedConsumer.Consume error", "error", err.Error())
			continue
		}
		var mission *models.MissionInfo
		err = json.Unmarshal(msg.Value, &mission)
		if err != nil {
			slog.Error("parse", "error", err)
			continue
		}
		if mission == nil || mission.ID == 0 {
			slog.Error("Invalid mission")
			continue
		}
		err = c.processor.ProcessMissionCreated(ctx, mission)
		if err != nil {
			slog.Error("ProcessMissionCreated error", "error", err.Error())
		}
	}
}
