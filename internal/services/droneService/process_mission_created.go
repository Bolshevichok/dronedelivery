package droneService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	"github.com/segmentio/kafka-go"
)

func (s *DroneService) ProcessMissionCreated(ctx context.Context, mission *models.MissionInfo) error {
	if mission == nil {
		return fmt.Errorf("mission is nil")
	}

	missionID := mission.ID
	if missionID == 0 {
		return fmt.Errorf("mission id is required")
	}

	// Fallback for demo: deterministic non-zero id.
	droneID := uint64(1)

	err := s.publishLifecycle(ctx, missionID, droneID, "assigned", "")
	if err != nil {
		return err
	}

	go s.simulateMission(ctx, missionID, droneID, mission)
	return nil
}

func (s *DroneService) simulateMission(ctx context.Context, missionID, droneID uint64, mission *models.MissionInfo) {
	time.Sleep(5 * time.Second)
	s.publishLifecycle(ctx, missionID, droneID, "picked_up", "")

	// Without launch base coordinates we simulate "no movement".
	startLat := mission.DestinationLat
	startLon := mission.DestinationLon

	go s.simulateTelemetry(ctx, missionID, droneID, startLat, startLon, mission.DestinationLat, mission.DestinationLon)

	time.Sleep(10 * time.Second)
	s.publishLifecycle(ctx, missionID, droneID, "delivered", "")
}

func (s *DroneService) publishLifecycle(ctx context.Context, missionID, droneID uint64, status, details string) error {
	_ = details
	event := &models.MissionLifecycleEvent{
		DroneID:   droneID,
		MissionID: missionID,
		Status:    status,
		Timestamp: time.Now().UTC(),
	}
	eventBytes, _ := json.Marshal(event)
	writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := s.lifecycleWriter.WriteMessages(writeCtx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", missionID)),
		Value: eventBytes,
	})
	if err != nil {
		log.Printf("Failed to publish lifecycle: %v", err)
		return err
	}
	return nil
}

func (s *DroneService) simulateTelemetry(ctx context.Context, missionID, droneID uint64, startLat, startLon, destLat, destLon float64) {
	steps := 10
	for i := 0; i < steps; i++ {
		lat := startLat + (destLat-startLat)*float64(i)/float64(steps)
		lon := startLon + (destLon-startLon)*float64(i)/float64(steps)
		telemetry := models.DroneTelemetry{
			DroneID:   droneID,
			MissionID: missionID,
			Lat:       lat,
			Lon:       lon,
			Alt:       100.0,
			Timestamp: time.Now(),
		}
		eventBytes, _ := json.Marshal(telemetry)
		writeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := s.telemetryWriter.WriteMessages(writeCtx, kafka.Message{
			Key:   []byte(fmt.Sprintf("%d", droneID)),
			Value: eventBytes,
		})
		cancel()
		if err != nil {
			log.Printf("Failed to publish telemetry: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}
