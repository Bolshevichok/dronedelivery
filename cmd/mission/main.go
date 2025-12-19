package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/bootstrap"
	missionv1 "github.com/Bolshevichok/dronedelivery/internal/pb/mission/v1"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	// Initialize dependencies (DB, Kafka, etc.)
	deps, err := bootstrap.InitMissionService(cfg)
	if err != nil {
		log.Fatalf("Failed to init dependencies: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register gRPC server
	missionv1.RegisterMissionServiceServer(grpcServer, missionService.NewMissionService(deps))

	// Start consumer for missions.lifecycle
	go func() {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)},
			Topic:    cfg.Kafka.MissionsLifecycleTopic,
			GroupID:  "mission-service",
			MinBytes: 10e3,
			MaxBytes: 10e6,
		})
		defer r.Close()

		log.Println("Mission service consuming missions.lifecycle")

		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading lifecycle: %v", err)
				continue
			}

			var event map[string]interface{}
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Error unmarshaling lifecycle: %v", err)
				continue
			}

			missionID := uint64(event["mission_id"].(float64))
			droneID := uint64(event["drone_id"].(float64))
			status := event["status"].(string)
			log.Printf("Updating mission %d status to %s by drone %d", missionID, status, droneID)

			// Update mission status in DB
			err = deps.Storage.UpdateMissionStatus(context.Background(), missionID, status)
			if err != nil {
				log.Printf("Failed to update mission status: %v", err)
				continue
			}

			// If assigned, insert mission_drone
			if status == "assigned" {
				missions, err := deps.Storage.GetMissionsByIDs(context.Background(), []uint64{missionID})
				if err != nil || len(missions) == 0 {
					log.Printf("Failed to load mission for assignment: %v", err)
					continue
				}
				mission := missions[0]

				missionDrone := &pgstorage.MissionDrone{
					MissionID:        missionID,
					DroneID:          droneID,
					AssignedBy:       mission.OperatorID,
					AssignedAt:       time.Now(),
					PlannedPayloadKg: mission.PayloadKg,
				}
				err = deps.Storage.UpsertMissionDrones(context.Background(), []*pgstorage.MissionDrone{missionDrone})
				if err != nil {
					log.Printf("Failed to assign drone: %v", err)
				}
			}
		}
	}()

	reflection.Register(grpcServer)

	// Start servers
	go func() {
		lis, err := net.Listen("tcp", ":8080") // gRPC port
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		log.Println("gRPC server listening on :8080")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	select {}
}
