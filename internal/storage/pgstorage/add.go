package pgstorage

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (storage *PGstorage) UpsertOperators(ctx context.Context, operators []*Operator) error {
	query := storage.upsertOperatorsQuery(operators)
	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}
	_, err = storage.db.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}
	return err
}

func (storage *PGstorage) upsertOperatorsQuery(operators []*Operator) squirrel.Sqlizer {
	q := squirrel.Insert(operatorTableName).Columns(OperatorEmailCol, OperatorNameCol, OperatorCreatedAtCol).
		PlaceholderFormat(squirrel.Dollar)
	for _, op := range operators {
		q = q.Values(op.Email, op.Name, op.CreatedAt)
	}
	return q
}

func (storage *PGstorage) UpsertLaunchBases(ctx context.Context, launchBases []*LaunchBase) error {
	query := storage.upsertLaunchBasesQuery(launchBases)
	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}
	_, err = storage.db.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}
	return err
}

func (storage *PGstorage) upsertLaunchBasesQuery(launchBases []*LaunchBase) squirrel.Sqlizer {
	q := squirrel.Insert(launchBaseTableName).Columns(LaunchBaseNameColumn, LaunchBaseLatColumn, LaunchBaseLonColumn, LaunchBaseAltColumn, LaunchBaseCreatedAtColumn).
		PlaceholderFormat(squirrel.Dollar)
	for _, lb := range launchBases {
		q = q.Values(lb.Name, lb.lat, lb.lon, lb.alt, lb.CreatedAt)
	}
	return q
}

func (storage *PGstorage) UpsertDrones(ctx context.Context, drones []*Drone) error {
	query := storage.upsertDronesQuery(drones)
	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}
	_, err = storage.db.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}
	return err
}

func (storage *PGstorage) upsertDronesQuery(drones []*Drone) squirrel.Sqlizer {
	q := squirrel.Insert(droneTableName).Columns(DroneSerialColumn, DroneModelColumn, DroneStatusColumn, DroneLaunchBaseIDColumn, DroneCreatedAtColumn).
		PlaceholderFormat(squirrel.Dollar)
	for _, d := range drones {
		q = q.Values(d.serial, d.model, d.status, d.launchbase, d.CreatedAt)
	}
	return q
}

func (storage *PGstorage) UpsertMissions(ctx context.Context, missions []*Mission) error {
	query := storage.upsertMissionsQuery(missions)
	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}
	_, err = storage.db.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}
	return err
}

func (storage *PGstorage) upsertMissionsQuery(missions []*Mission) squirrel.Sqlizer {
	q := squirrel.Insert(missionTableName).Columns(MissionOperatorIDColumn, MissionLaunchBaseIDColumn, MissionStatusColumn, MissionDestinationLatColumn, MissionDestinationLonColumn, MissionDestinationAltColumn, MissionPayloadKgColumn, MissionCreatedAtColumn).
		PlaceholderFormat(squirrel.Dollar)
	for _, m := range missions {
		q = q.Values(m.operator_id, m.launch_base_id, m.status, m.destination_lat, m.destination_lon, m.destination_alt, m.payload_kg, m.created_at)
	}
	return q
}

func (storage *PGstorage) UpsertMissionDrones(ctx context.Context, missionDrones []*MissionDrone) error {
	query := storage.upsertMissionDronesQuery(missionDrones)
	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}
	_, err = storage.db.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}
	return err
}

func (storage *PGstorage) upsertMissionDronesQuery(missionDrones []*MissionDrone) squirrel.Sqlizer {
	q := squirrel.Insert(missionDroneTableName).Columns(MissionDroneMissionIDColumn, MissionDroneDroneIDColumn, MissionDroneAssignedByColumn, MissionDroneAssignedAtColumn, MissionDronePlannedPayloadKgColumn).
		PlaceholderFormat(squirrel.Dollar)
	for _, md := range missionDrones {
		q = q.Values(md.MissionID, md.DroneID, md.assigned_by, md.assigned_at, md.planned_payload_kg)
	}
	return q
}
