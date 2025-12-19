package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/go-redis/redis/v8"
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

	// Init Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})
	defer rdb.Close()

	// Init consumer for drone.telemetry
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)},
		Topic:    cfg.Kafka.DroneTelemetryTopic,
		GroupID:  "telemetry-extractor",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Println("Telemetry extractor started, consuming drone.telemetry")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var telemetry map[string]interface{}
		if err := json.Unmarshal(m.Value, &telemetry); err != nil {
			log.Printf("Error unmarshaling telemetry: %v", err)
			continue
		}

		droneID := uint64(telemetry["drone_id"].(float64))
		key := fmt.Sprintf("telemetry:%d", droneID)

		// Храним последнюю телеметрию по дрону (TTL 1 час).
		err = rdb.Set(context.Background(), key, string(m.Value), time.Hour).Err()
		if err != nil {
			log.Printf("Failed to store telemetry in Redis: %v", err)
		}
	}
}
