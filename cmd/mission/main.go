package main

import (
	"log"
	"net"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/bootstrap"
	missionv1 "github.com/Bolshevichok/dronedelivery/internal/pb/mission/v1"
	"github.com/Bolshevichok/dronedelivery/internal/services/missionService"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize dependencies (DB, Kafka, etc.)
	deps, err := bootstrap.InitMissionService(cfg)
	if err != nil {
		log.Fatalf("Failed to init dependencies: %v", err)
	}

	// gRPC server
	grpcServer := grpc.NewServer()
	missionv1.RegisterMissionServiceServer(grpcServer, missionService.NewMissionService(deps))
	reflection.Register(grpcServer)

	// gRPC-Gateway (HTTP)
	// mux := runtime.NewServeMux()
	// err = missionv1.RegisterMissionServiceHandlerFromEndpoint(context.Background(), mux, "localhost:8080", []grpc.DialOption{grpc.WithInsecure()})
	// if err != nil {
	// 	log.Fatalf("Failed to register gateway: %v", err)
	// }

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

	// log.Println("HTTP gateway listening on :8081")
	// if err := http.ListenAndServe(":8081", mux); err != nil {
	// 	log.Fatalf("Failed to serve HTTP: %v", err)
	// }
}
