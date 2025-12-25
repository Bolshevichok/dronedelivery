package bootstrap

import (
	"context"
	"log"
	"net"

	"github.com/Bolshevichok/dronedelivery/internal/api/mission_api"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/mission_created_consumer"
	"github.com/Bolshevichok/dronedelivery/internal/consumer/telemetry_consumer"
	pb_mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func AppRun(api *mission_api.MissionAPI, consumer *mission_created_consumer.MissionCreatedConsumerImpl) {
	go consumer.Start(context.Background())

	grpcServer := grpc.NewServer()
	pb_mission_api.RegisterMissionServiceServer(grpcServer, api)
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

func AppRunConsumer(consumer *mission_created_consumer.MissionCreatedConsumerImpl) {
	consumer.Start(context.Background())
}

func AppRunTelemetryConsumer(consumer *telemetry_consumer.TelemetryConsumerImpl) {
	consumer.Start(context.Background())
}
