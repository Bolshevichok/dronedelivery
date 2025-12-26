package missionService

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MissionService) UpsertMissions(ctx context.Context, missions []*models.Mission) ([]*models.Mission, error) {
	for _, mission := range missions {
		upsertedMissions, err := s.missionStorage.UpsertMissions(ctx, []*models.Mission{mission})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to upsert mission: %v", err)
		}
		mission = upsertedMissions[0]

		missionInfo := &models.MissionInfo{
			ID:             mission.ID,
			Status:         mission.Status,
			DestinationLat: mission.DestinationLat,
			DestinationLon: mission.DestinationLon,
			DestinationAlt: mission.DestinationAlt,
			PayloadKg:      mission.PayloadKg,
		}
		eventBytes, _ := json.Marshal(missionInfo)
		err = s.kafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(fmt.Sprintf("%d", mission.ID)),
			Value: eventBytes,
		})
		if err != nil {
			slog.Error("не удалось отправить missions.created", "mission_id", mission.ID, "err", err)
		}
	}

	return missions, nil
}
