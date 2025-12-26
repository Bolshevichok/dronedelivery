package bootstrap

import (
	"context"
	"log"
	"net"

	"github.com/Bolshevichok/dronedelivery/internal/api/mission_api"
	pb_mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Consumer interface {
	Consume(ctx context.Context)
}

func AppRun(api *mission_api.MissionAPI, consumers ...Consumer) {
	for _, consumer := range consumers {
		go consumer.Consume(context.Background())
	}

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

func AppRunConsumers(consumers ...Consumer) {
	for _, consumer := range consumers {
		go consumer.Consume(context.Background())
	}
	// Block forever, since no server to run
	select {}
}
