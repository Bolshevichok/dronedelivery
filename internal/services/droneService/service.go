package droneService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
	"github.com/segmentio/kafka-go"
)

type DroneService struct {
	storage         *pgstorage.PGstorage
	lifecycleWriter *kafka.Writer
	telemetryWriter *kafka.Writer
}

func NewDroneService(storage *pgstorage.PGstorage, lifecycleWriter, telemetryWriter *kafka.Writer) *DroneService {
	return &DroneService{
		storage:         storage,
		lifecycleWriter: lifecycleWriter,
		telemetryWriter: telemetryWriter,
	}
}

func (s *DroneService) ProcessMissionCreated(ctx context.Context, missionID uint64) {
	// Get mission
	missions, err := s.storage.GetMissionsByIDs(ctx, []uint64{missionID})
	if err != nil || len(missions) == 0 {
		log.Printf("Failed to get mission %d: %v", missionID, err)
		return
	}
	mission := missions[0]

	// Берём любой доступный дрон на нужной базе.
	drones, err := s.storage.GetAvailableDrones(ctx, mission.LaunchBaseID)
	if err != nil || len(drones) == 0 {
		log.Printf("No available drones for mission %d", missionID)
		return
	}
	drone := drones[0] // Для учебного стенда берём первый подходящий.

	// Publish lifecycle assigned
	s.publishLifecycle(ctx, missionID, drone.ID, "assigned", "")

	// Start simulation
	go s.simulateMission(ctx, missionID, drone.ID, mission)
}

func (s *DroneService) simulateMission(ctx context.Context, missionID, droneID uint64, mission *pgstorage.Mission) {
	// Simulate flight
	time.Sleep(5 * time.Second) // Picked up
	s.publishLifecycle(ctx, missionID, droneID, "picked_up", "")

	// Simulate telemetry during flight
	go s.simulateTelemetry(ctx, missionID, droneID, mission.DestinationLat, mission.DestinationLon)

	time.Sleep(10 * time.Second) // Delivered
	s.publishLifecycle(ctx, missionID, droneID, "delivered", "")
}

func (s *DroneService) publishLifecycle(ctx context.Context, missionID, droneID uint64, status, details string) {
	event := map[string]interface{}{
		"event_id":   fmt.Sprintf("event-%d-%s", missionID, status),
		"mission_id": missionID,
		"drone_id":   droneID,
		"status":     status,
		"details":    details,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	eventBytes, _ := json.Marshal(event)
	err := s.lifecycleWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", missionID)),
		Value: eventBytes,
	})
	if err != nil {
		log.Printf("Failed to publish lifecycle: %v", err)
	}
}

func (s *DroneService) simulateTelemetry(ctx context.Context, missionID, droneID uint64, destLat, destLon float64) {
	// Simulate moving from base to destination
	startLat, startLon := 55.7558, 37.6173 // Moscow base
	steps := 10
	for i := 0; i < steps; i++ {
		lat := startLat + (destLat-startLat)*float64(i)/float64(steps)
		lon := startLon + (destLon-startLon)*float64(i)/float64(steps)
		telemetry := map[string]interface{}{
			"drone_id":   droneID,
			"mission_id": missionID,
			"lat":        lat,
			"lon":        lon,
			"alt":        100.0,
			"timestamp":  time.Now().Format(time.RFC3339),
		}
		eventBytes, _ := json.Marshal(telemetry)
		err := s.telemetryWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(fmt.Sprintf("%d", droneID)),
			Value: eventBytes,
		})
		if err != nil {
			log.Printf("Failed to publish telemetry: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}
