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
	ID              uint64  `db:"id"`
	operator_id     uint64  `db:"operator_id"`
	launch_base_id  uint64  `db:"launch_base_id"`
	status          string  `db:"status"`
	destination_lat float64 `db:"destination_lat"`
	destination_lon float64 `db:"destination_lon"`
	destination_alt float64 `db:"destination_alt"`
	payload_kg      float64 `db:"payload"`
	created_at      string  `db:"created_at"`
	Operator        Operator
	LaunchBase      LaunchBase
	Drones          []Drone
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
	lat       float64 `db:"lat"`
	lon       float64 `db:"lon"`
	alt       float64 `db:"alt"`
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
	ID         uint64 `db:"id"`
	serial     string `db:"serial"`
	model      string `db:"model"`
	status     string `db:"status"`
	launchbase uint64 `db:"launch_base_id"`
	CreatedAt  string `db:"created_at"`
	LaunchBase LaunchBase
	Missions   []Mission
}

const (
	droneTableName = "drones"

	DroneIDColumn        = "id"
	DroneSerialColumn    = "serial"
	DroneModelColumn     = "model"
	DroneStatusColumn    = "status"
	DroneLaunchBaseIDCol = "launch_base_id"
	DroneCreatedAtColumn = "created_at"
)

type MissionDrone struct {
	MissionID          uint64  `db:"mission_id"`
	DroneID            uint64  `db:"drone_id"`
	assigned_by        uint64  `db:"assigned_by"`
	assigned_at        string  `db:"assigned_at"`
	planned_payload_kg float64 `db:"planned_payload_kg"`
	Mission            Mission
	Drone              Drone
	AssignedBy         Operator
}

const (
	missionDroneTableName = "mission_drones"

	MissionDroneMissionIDColumn        = "mission_id"
	MissionDroneDroneIDColumn          = "drone_id"
	MissionDroneAssignedByColumn       = "assigned_by"
	MissionDroneAssignedAtColumn       = "assigned_at"
	MissionDronePlannedPayloadKgColumn = "planned_payload_kg"
)
