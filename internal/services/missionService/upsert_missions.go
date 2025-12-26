package missionService

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MissionService) UpsertMissions(ctx context.Context, req *mission_api.UpsertMissionsRequest) (*mission_api.UpsertMissionsResponse, error) {
	for _, missionPb := range req.Missions {
		mission := &models.Mission{
			OperatorID:     missionPb.OpId,
			LaunchBaseID:   missionPb.BaseId,
			Status:         missionPb.Status,
			DestinationLat: missionPb.Lat,
			DestinationLon: missionPb.Lon,
			DestinationAlt: missionPb.Alt,
			PayloadKg:      missionPb.Payload,
			CreatedAt:      time.Now(),
		}

		missions, err := s.missionStorage.UpsertMissions(ctx, []*models.Mission{mission})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to upsert mission: %v", err)
		}
		mission = missions[0]

		event := map[string]interface{}{
			"event_id":   fmt.Sprintf("mission-created-%d", mission.ID),
			"mission_id": mission.ID,
			"base_id":    mission.LaunchBaseID,
			"payload":    missionPb,
			"timestamp":  time.Now().Format(time.RFC3339),
		}
		eventBytes, _ := json.Marshal(event)
		err = s.kafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(fmt.Sprintf("%d", mission.ID)),
			Value: eventBytes,
		})
		if err != nil {
			slog.Error("не удалось отправить missions.created", "mission_id", mission.ID, "err", err)
		}
	}

	return &mission_api.UpsertMissionsResponse{}, nil
}
