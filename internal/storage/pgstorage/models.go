package pgstorage

import (
	"github.com/Bolshevichok/dronedelivery/internal/models"
)

type Operator = models.Operator

const (
	operatorTableName = "operators"

	OperatorIDColumn     = "id"
	OperatorEmailCol     = "email"
	OperatorNameCol      = "name"
	OperatorCreatedAtCol = "created_at"
)

type Mission = models.Mission

const (
	missionTableName = "missions"

	MissionIDColumn             = "id"
	MissionOperatorIDColumn     = "operator_id"
	MissionLaunchBaseIDColumn   = "launch_base_id"
	MissionStatusColumn         = "status"
	MissionDestinationLatColumn = "destination_lat"
	MissionDestinationLonColumn = "destination_lon"
	MissionDestinationAltColumn = "destination_alt"
	MissionPayloadKgColumn      = "payload"
	MissionCreatedAtColumn      = "created_at"
)

type LaunchBase = models.LaunchBase

const (
	launchBaseTableName = "launch_bases"

	LaunchBaseIDColumn        = "id"
	LaunchBaseNameColumn      = "name"
	LaunchBaseLatColumn       = "lat"
	LaunchBaseLonColumn       = "lon"
	LaunchBaseAltColumn       = "alt"
	LaunchBaseCreatedAtColumn = "created_at"
)

type Drone = models.Drone

const (
	droneTableName = "drones"

	DroneIDColumn           = "id"
	DroneSerialColumn       = "serial"
	DroneModelColumn        = "model"
	DroneStatusColumn       = "status"
	DroneLaunchBaseIDColumn = "launch_base_id"
	DroneCreatedAtColumn    = "created_at"
)

type MissionDrone = models.MissionDrone

const (
	missionDroneTableName = "mission_drones"

	MissionDroneMissionIDColumn        = "mission_id"
	MissionDroneDroneIDColumn          = "drone_id"
	MissionDroneAssignedByColumn       = "assigned_by"
	MissionDroneAssignedAtColumn       = "assigned_at"
	MissionDronePlannedPayloadKgColumn = "planned_payload_kg"
)
