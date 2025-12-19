package pgstorage

type Operator struct {
	ID        uint64 `db:"id"`
	Email     string `db:"email"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	Missions  []Mission
}

const (
	operatorTableName = "operators"

	OperatorIDColumn     = "id"
	OperatorEmailCol     = "email"
	OperatorNameCol      = "name"
	OperatorCreatedAtCol = "created_at"
)

type Mission struct {
	ID             uint64  `db:"id"`
	OperatorID     uint64  `db:"operator_id"`
	LaunchBaseID   uint64  `db:"launch_base_id"`
	Status         string  `db:"status"`
	DestinationLat float64 `db:"destination_lat"`
	DestinationLon float64 `db:"destination_lon"`
	DestinationAlt float64 `db:"destination_alt"`
	PayloadKg      float64 `db:"payload"`
	CreatedAt      string  `db:"created_at"`
	Operator       Operator
	LaunchBase     LaunchBase
	Drones         []Drone
}

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

type LaunchBase struct {
	ID        uint64  `db:"id"`
	Name      string  `db:"name"`
	Lat       float64 `db:"lat"`
	Lon       float64 `db:"lon"`
	Alt       float64 `db:"alt"`
	CreatedAt string  `db:"created_at"`
	Missions  []Mission
	Drones    []Drone
}

const (
	launchBaseTableName = "launch_bases"

	LaunchBaseIDColumn        = "id"
	LaunchBaseNameColumn      = "name"
	LaunchBaseLatColumn       = "lat"
	LaunchBaseLonColumn       = "lon"
	LaunchBaseAltColumn       = "alt"
	LaunchBaseCreatedAtColumn = "created_at"
)

type Drone struct {
	ID           uint64 `db:"id"`
	Serial       string `db:"serial"`
	Model        string `db:"model"`
	Status       string `db:"status"`
	LaunchBaseID uint64 `db:"launch_base_id"`
	CreatedAt    string `db:"created_at"`
	LaunchBase   LaunchBase
	Missions     []Mission
}

const (
	droneTableName = "drones"

	DroneIDColumn           = "id"
	DroneSerialColumn       = "serial"
	DroneModelColumn        = "model"
	DroneStatusColumn       = "status"
	DroneLaunchBaseIDColumn = "launch_base_id"
	DroneCreatedAtColumn    = "created_at"
)

type MissionDrone struct {
	MissionID          uint64  `db:"mission_id"`
	DroneID            uint64  `db:"drone_id"`
	AssignedBy         uint64  `db:"assigned_by"`
	AssignedAt         string  `db:"assigned_at"`
	PlannedPayloadKg   float64 `db:"planned_payload_kg"`
	Mission            Mission
	Drone              Drone
	AssignedByOperator Operator
}

const (
	missionDroneTableName = "mission_drones"

	MissionDroneMissionIDColumn        = "mission_id"
	MissionDroneDroneIDColumn          = "drone_id"
	MissionDroneAssignedByColumn       = "assigned_by"
	MissionDroneAssignedAtColumn       = "assigned_at"
	MissionDronePlannedPayloadKgColumn = "planned_payload_kg"
)
