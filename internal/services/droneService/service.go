package droneService

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal"
	"github.com/Bolshevichok/dronedelivery/internal/models"
	"github.com/segmentio/kafka-go"
)

// DroneService interface
type DroneService interface {
	ProcessMissionCreated(ctx context.Context, missionID uint64)
}

type DroneServiceImpl struct {
	storage         internal.Storage
	lifecycleWriter *kafka.Writer
	telemetryWriter *kafka.Writer
}

func NewDroneService(storage internal.Storage, cfg *config.Config) DroneService {
	lifecycleWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.DroneLifecycleTopic,
		Balancer: &kafka.LeastBytes{},
	}
	telemetryWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.DroneTelemetryTopic,
		Balancer: &kafka.LeastBytes{},
	}

	return &DroneServiceImpl{
		storage:         storage,
		lifecycleWriter: lifecycleWriter,
		telemetryWriter: telemetryWriter,
	}
}

func (s *DroneServiceImpl) ProcessMissionCreated(ctx context.Context, missionID uint64) {
	missions, err := s.storage.GetMissionsByIDs(ctx, []uint64{missionID})
	if err != nil || len(missions) == 0 {
		log.Printf("Failed to get mission %d: %v", missionID, err)
		return
	}
	mission := missions[0]

	drones, err := s.storage.GetAvailableDrones(ctx, mission.LaunchBaseID)
	if err != nil || len(drones) == 0 {
		log.Printf("No available drones for mission %d", missionID)
		return
	}
	drone := drones[0]

	s.publishLifecycle(ctx, missionID, drone.ID, "assigned", "")

	go s.simulateMission(ctx, missionID, drone.ID, mission)
}

func (s *DroneServiceImpl) simulateMission(ctx context.Context, missionID, droneID uint64, mission *models.Mission) {
	time.Sleep(5 * time.Second)
	s.publishLifecycle(ctx, missionID, droneID, "picked_up", "")

	launchBases, err := s.storage.GetLaunchBasesByIDs(ctx, []uint64{mission.LaunchBaseID})
	if err != nil || len(launchBases) == 0 {
		log.Printf("Failed to get launch base %d: %v", mission.LaunchBaseID, err)
		return
	}
	launchBase := launchBases[0]

	go s.simulateTelemetry(ctx, missionID, droneID, launchBase.Lat, launchBase.Lon, mission.DestinationLat, mission.DestinationLon)

	time.Sleep(10 * time.Second)
	s.publishLifecycle(ctx, missionID, droneID, "delivered", "")
}

func (s *DroneServiceImpl) publishLifecycle(ctx context.Context, missionID, droneID uint64, status, details string) {
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

func (s *DroneServiceImpl) simulateTelemetry(ctx context.Context, missionID, droneID uint64, startLat, startLon, destLat, destLon float64) {
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
