package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Bolshevichok/dronedelivery/config"
	missionv1 "github.com/Bolshevichok/dronedelivery/internal/pb/mission/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cfg        *config.Config
	grpcClient missionv1.MissionServiceClient
)

func main() {
	cfgPath := os.Getenv("configPath")
	if cfgPath == "" {
		cfgPath = "config.yaml"
	}

	var err error
	cfg, err = config.LoadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to gRPC
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.MissionService.Host, cfg.MissionService.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	grpcClient = missionv1.NewMissionServiceClient(conn)

	var rootCmd = &cobra.Command{Use: "cli"}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new mission",
		Run:   createMission,
	}
	createCmd.Flags().Uint64("operator-id", 1, "Operator ID")
	createCmd.Flags().Uint64("base-id", 1, "Launch base ID")
	createCmd.Flags().Float64("lat", 0, "Destination latitude")
	createCmd.Flags().Float64("lon", 0, "Destination longitude")
	createCmd.Flags().Float64("alt", 0, "Destination altitude")
	createCmd.Flags().Float64("payload", 0, "Payload kg")

	var getCmd = &cobra.Command{
		Use:   "get [mission-id]",
		Short: "Get mission by ID",
		Args:  cobra.ExactArgs(1),
		Run:   getMission,
	}

	var watchCmd = &cobra.Command{
		Use:   "watch [mission-id]",
		Short: "Watch mission updates",
		Args:  cobra.ExactArgs(1),
		Run:   watchMission,
	}

	rootCmd.AddCommand(createCmd, getCmd, watchCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func createMission(cmd *cobra.Command, args []string) {
	operatorID, _ := cmd.Flags().GetUint64("operator-id")
	baseID, _ := cmd.Flags().GetUint64("base-id")
	lat, _ := cmd.Flags().GetFloat64("lat")
	lon, _ := cmd.Flags().GetFloat64("lon")
	alt, _ := cmd.Flags().GetFloat64("alt")
	payload, _ := cmd.Flags().GetFloat64("payload")

	req := &missionv1.CreateMissionRequest{
		OperatorId:     operatorID,
		LaunchBaseId:   baseID,
		DestinationLat: lat,
		DestinationLon: lon,
		DestinationAlt: alt,
		PayloadKg:      payload,
	}

	resp, err := grpcClient.CreateMission(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to create mission: %v", err)
	}

	fmt.Printf("Mission created with ID: %d\n", resp.MissionId)
}

func getMission(cmd *cobra.Command, args []string) {
	missionID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		log.Fatalf("Invalid mission ID: %v", err)
	}

	req := &missionv1.GetMissionRequest{MissionId: missionID}

	resp, err := grpcClient.GetMission(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get mission: %v", err)
	}

	mission := resp.Mission
	fmt.Printf("Mission ID: %d\n", mission.Id)
	fmt.Printf("Status: %s\n", mission.Status)
	fmt.Printf("Destination: %.4f, %.4f, %.2f\n", mission.DestinationLat, mission.DestinationLon, mission.DestinationAlt)
	fmt.Printf("Payload: %.2f kg\n", mission.PayloadKg)
	fmt.Printf("Created: %s\n", mission.CreatedAt)
}

func watchMission(cmd *cobra.Command, args []string) {
	missionID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		log.Fatalf("Invalid mission ID: %v", err)
	}

	req := &missionv1.WatchMissionRequest{MissionId: missionID}

	stream, err := grpcClient.WatchMission(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to watch mission: %v", err)
	}

	for {
		update, err := stream.Recv()
		if err != nil {
			log.Fatalf("Stream error: %v", err)
		}

		fmt.Printf("[%s] Status: %s, Telemetry: lat=%.4f lon=%.4f alt=%.2f\n",
			update.Timestamp, update.Status, update.Telemetry.Lat, update.Telemetry.Lon, update.Telemetry.Alt)
	}
}
