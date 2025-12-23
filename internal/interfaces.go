package internal

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
)

type Storage interface {
	UpsertMissions(ctx context.Context, missions []*models.Mission) ([]*models.Mission, error)
	GetMissionsByIDs(ctx context.Context, IDs []uint64) ([]*models.Mission, error)
	UpsertOperators(ctx context.Context, operators []*models.Operator) error
	GetOperatorsByIDs(ctx context.Context, IDs []uint64) ([]*models.Operator, error)
	UpsertLaunchBases(ctx context.Context, launchBases []*models.LaunchBase) error
	GetLaunchBasesByIDs(ctx context.Context, IDs []uint64) ([]*models.LaunchBase, error)
	UpsertDrones(ctx context.Context, drones []*models.Drone) error
	GetDronesByIDs(ctx context.Context, IDs []uint64) ([]*models.Drone, error)
	UpsertMissionDrones(ctx context.Context, missionDrones []*models.MissionDrone) error
	GetMissionDronesByMissionIDs(ctx context.Context, missionIDs []uint64) ([]*models.MissionDrone, error)
	GetAvailableDrones(ctx context.Context, launchBaseID uint64) ([]*models.Drone, error)
	GetMissionDronesByMissionID(ctx context.Context, missionID uint64) ([]*models.MissionDrone, error)
	UpdateMissionStatus(ctx context.Context, missionID uint64, status string) error
}
