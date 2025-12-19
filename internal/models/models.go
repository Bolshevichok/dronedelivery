package models

type Operator struct {
	ID        uint64
	Email     string
	Name      string
	CreatedAt string
	Missions  []Mission
}

type Mission struct {
	ID             uint64
	OperatorID     uint64
	LaunchBaseID   uint64
	Status         string
	DestinationLat float64
	DestinationLon float64
	DestinationAlt float64
	PayloadKg      float64
	CreatedAt      string
	Operator       Operator
	LaunchBase     LaunchBase
	Drones         []Drone
}

type LaunchBase struct {
	ID        uint64
	Name      string
	Lat       float64
	Lon       float64
	Alt       float64
	CreatedAt string
	Missions  []Mission
	Drones    []Drone
}

type Drone struct {
	ID           uint64
	Serial       string
	Model        string
	Status       string
	LaunchBaseID uint64
	CreatedAt    string
	LaunchBase   LaunchBase
	Missions     []Mission
}

type MissionDrone struct {
	MissionID          uint64
	DroneID            uint64
	AssignedBy         uint64
	AssignedAt         string
	PlannedPayloadKg   float64
	Mission            Mission
	Drone              Drone
	AssignedByOperator Operator
}
