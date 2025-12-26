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

func (s *DroneService) ProcessMissionCreated(ctx context.Context, missionID uint64) error {
	missions, err := s.droneStorage.GetMissionsByIDs(ctx, []uint64{missionID})
	if err != nil {
		log.Printf("Failed to get mission %d: %v", missionID, err)
		return err
	}
	if len(missions) == 0 {
		log.Printf("Mission %d not found", missionID)
		return fmt.Errorf("mission not found")
	}
	mission := missions[0]

	drones, err := s.droneStorage.GetAvailableDrones(ctx, mission.LaunchBaseID)
	if err != nil {
		log.Printf("Failed to get available drones for mission %d: %v", missionID, err)
		return err
	}
	if len(drones) == 0 {
		log.Printf("No available drones for mission %d", missionID)
		return fmt.Errorf("no available drones")
	}
	drone := drones[0]

	err = s.publishLifecycle(ctx, missionID, drone.ID, "assigned", "")
	if err != nil {
		return err
	}

	go s.simulateMission(ctx, missionID, drone.ID, mission)
	return nil
}

func (s *DroneService) simulateMission(ctx context.Context, missionID, droneID uint64, mission *models.Mission) {
	time.Sleep(5 * time.Second)
	s.publishLifecycle(ctx, missionID, droneID, "picked_up", "")

	launchBases, err := s.droneStorage.GetLaunchBasesByIDs(ctx, []uint64{mission.LaunchBaseID})
	if err != nil || len(launchBases) == 0 {
		log.Printf("Failed to get launch base %d: %v", mission.LaunchBaseID, err)
		return
	}
	launchBase := launchBases[0]

	go s.simulateTelemetry(ctx, missionID, droneID, launchBase.Lat, launchBase.Lon, mission.DestinationLat, mission.DestinationLon)

	time.Sleep(10 * time.Second)
	s.publishLifecycle(ctx, missionID, droneID, "delivered", "")
}

func (s *DroneService) publishLifecycle(ctx context.Context, missionID, droneID uint64, status, details string) error {
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
		return err
	}
	return nil
}

func (s *DroneService) simulateTelemetry(ctx context.Context, missionID, droneID uint64, startLat, startLon, destLat, destLon float64) {
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
