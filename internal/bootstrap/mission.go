package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/api/mission_api"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/mission_created_consumer"
	missionv1 "github.com/Bolshevichok/dronedelivery/internal/pb/mission/v1"
	"github.com/Bolshevichok/dronedelivery/internal/services/droneService"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
	"github.com/Bolshevichok/dronedelivery/internal/services/processors/mission_processor"
	"github.com/Bolshevichok/dronedelivery/internal/storage/pgstorage"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MissionComponents struct {
	MissionService         missionService.MissionService
	MissionProcessor       mission_processor.MissionProcessor
	MissionCreatedConsumer mission_created_consumer.MissionCreatedConsumer
	MissionAPI             *mission_api.MissionAPI
}

func InitMissionComponents(cfg *config.Config) (*MissionComponents, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		fmt.Sprintf("%d", cfg.Database.Port),
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	storage, err := pgstorage.NewPGStorge(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to init storage: %w", err)
	}

	// Kafka writers
	createdWriter := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
		Topic:    cfg.Kafka.MissionsCreatedTopic,
		Balancer: &kafka.LeastBytes{},
	}

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

	// Kafka reader for consumer
	createdReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)},
		Topic:   cfg.Kafka.MissionsCreatedTopic,
		GroupID: "mission-created-consumer",
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	})

	// Initialize services
	missionSvcDeps := &missionService.Dependencies{
		Storage:     storage,
		KafkaWriter: createdWriter,
		RedisClient: redisClient,
	}
	missionSvc := missionService.NewMissionService(missionSvcDeps)

	droneSvc := droneService.NewDroneService(storage, lifecycleWriter, telemetryWriter)

	// Initialize processor
	processorDeps := &mission_processor.Dependencies{
		Storage:        storage,
		MissionService: missionSvc,
	}
	processor := mission_processor.NewMissionProcessor(processorDeps)

	// Initialize consumer
	consumerDeps := &mission_created_consumer.Dependencies{
		DroneService: droneSvc,
		KafkaReader:  createdReader,
	}
	consumer := mission_created_consumer.NewMissionCreatedConsumer(consumerDeps)

	// Initialize API
	missionAPI := mission_api.NewMissionAPI(missionSvc)

	return &MissionComponents{
		MissionService:         missionSvc,
		MissionProcessor:       processor,
		MissionCreatedConsumer: consumer,
		MissionAPI:             missionAPI,
	}, nil
}

func AppRun(api *mission_api.MissionAPI, consumer mission_created_consumer.MissionCreatedConsumer) {
	// Start the consumer
	go func() {
		if err := consumer.Start(context.Background()); err != nil {
			log.Fatalf("Failed to start consumer: %v", err)
		}
	}()

	grpcServer := grpc.NewServer()
	missionv1.RegisterMissionServiceServer(grpcServer, api)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("gRPC server listening on :8080")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
