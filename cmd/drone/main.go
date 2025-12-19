package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/services/droneService"
	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
	"github.com/segmentio/kafka-go"
)

func main() {
	cfgPath := os.Getenv("configPath")
	if cfgPath == "" {
		cfgPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Init storage
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)
	storage, err := pgstorage.NewPGStorge(connString)
	if err != nil {
		log.Fatalf("Failed to init storage: %v", err)
	}

	// Init Kafka writers
	lifecycleWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.MissionsLifecycleTopic,
		Balancer: &kafka.LeastBytes{},
	}
	telemetryWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.DroneTelemetryTopic,
		Balancer: &kafka.LeastBytes{},
	}

	// Init service
	droneSvc := droneService.NewDroneService(storage, lifecycleWriter, telemetryWriter)

	// Init consumer for missions.created
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)},
		Topic:    cfg.Kafka.MissionsCreatedTopic,
		GroupID:  "drone-service",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	log.Println("Drone service started, consuming missions.created")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var event map[string]interface{}
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			continue
		}

		missionID := uint64(event["mission_id"].(float64))
		log.Printf("Received mission created: %d", missionID)

		// Process mission
		go droneSvc.ProcessMissionCreated(context.Background(), missionID)
	}
}
