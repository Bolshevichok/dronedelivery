package mission_created_consumer

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

type missionProcessor interface {
	ProcessMissionCreated(ctx context.Context, mission *models.Mission) error
}

type MissionCreatedConsumer interface {
	Consume(ctx context.Context)
}

type MissionCreatedConsumerImpl struct {
	processor   missionProcessor
	kafkaBroker []string
	topicName   string
}

func NewMissionCreatedConsumer(processor missionProcessor, kafkaBroker []string, topicName string) *MissionCreatedConsumerImpl {
	return &MissionCreatedConsumerImpl{
		processor:   processor,
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
	}
}
