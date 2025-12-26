package mission_lifecycle_consumer

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

type missionLifecycleProcessor interface {
	ProcessMissionLifecycle(ctx context.Context, mission *models.Mission) error
}

type MissionLifecycleConsumer interface {
	Consume(ctx context.Context)
}

type MissionLifecycleConsumerImpl struct {
	processor   missionLifecycleProcessor
	kafkaBroker []string
	topicName   string
}

func NewMissionLifecycleConsumer(processor missionLifecycleProcessor, kafkaBroker []string, topicName string) *MissionLifecycleConsumerImpl {
	return &MissionLifecycleConsumerImpl{
		processor:   processor,
		kafkaBroker: kafkaBroker,
		topicName:   topicName,
	}
}
