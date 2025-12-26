package mission_created_consumer

import (
	"context"
)

type missionProcessor interface {
	ProcessMissionCreated(ctx context.Context, missionID uint64)
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
