package models

import (
	"encoding/json"
	"time"
)

type Operator struct {
	ID        uint64
	Email     string
	Name      string
	CreatedAt time.Time
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
	CreatedAt      time.Time
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
	CreatedAt time.Time
	Missions  []Mission
	Drones    []Drone
}

type Drone struct {
	ID           uint64
	Serial       string
	Model        string
	Status       string
	LaunchBaseID uint64
	CreatedAt    time.Time
	LaunchBase   LaunchBase
	Missions     []Mission
}

type MissionInfo struct {
	ID             uint64
	Status         string
	DestinationLat float64
	DestinationLon float64
	DestinationAlt float64
	PayloadKg      float64
}

type MissionDrone struct {
	MissionID          uint64
	DroneID            uint64
	AssignedBy         uint64
	AssignedAt         time.Time
	PlannedPayloadKg   float64
	Mission            Mission
	Drone              Drone
	AssignedByOperator Operator
}

type MissionLifecycleEvent struct {
	DroneID   uint64
	MissionID uint64
	Status    string
	Timestamp time.Time
}

func (e MissionLifecycleEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"drone_id":   e.DroneID,
		"mission_id": e.MissionID,
		"status":     e.Status,
		"timestamp":  e.Timestamp,
	})
}

func (e *MissionLifecycleEvent) UnmarshalJSON(data []byte) error {
	if e == nil {
		return nil
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if v, ok := raw["drone_id"]; ok {
		_ = json.Unmarshal(v, &e.DroneID)
	}
	if v, ok := raw["mission_id"]; ok {
		_ = json.Unmarshal(v, &e.MissionID)
	}
	if v, ok := raw["status"]; ok {
		_ = json.Unmarshal(v, &e.Status)
	}
	if v, ok := raw["timestamp"]; ok {
		_ = json.Unmarshal(v, &e.Timestamp)
	}
	return nil
}

type DroneTelemetry struct {
	DroneID   uint64
	MissionID uint64
	Lat       float64
	Lon       float64
	Alt       float64
	Timestamp time.Time
}
