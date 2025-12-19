package pgstorage

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (storage *PGstorage) GetOperatorsByIDs(ctx context.Context, IDs []uint64) ([]*Operator, error) {
	query := storage.getOperatorsQuery(IDs)
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()
	var operators []*Operator
	for rows.Next() {
		var op Operator
		if err := rows.Scan(&op.ID, &op.Email, &op.Name, &op.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		operators = append(operators, &op)
	}
	return operators, nil
}

func (storage *PGstorage) getOperatorsQuery(IDs []uint64) squirrel.Sqlizer {
	q := squirrel.Select(OperatorIDColumn, OperatorEmailCol, OperatorNameCol, OperatorCreatedAtCol).From(operatorTableName).
		Where(squirrel.Eq{OperatorIDColumn: IDs}).PlaceholderFormat(squirrel.Dollar)
	return q
}

func (storage *PGstorage) GetLaunchBasesByIDs(ctx context.Context, IDs []uint64) ([]*LaunchBase, error) {
	query := storage.getLaunchBasesQuery(IDs)
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()
	var launchBases []*LaunchBase
	for rows.Next() {
		var lb LaunchBase
		if err := rows.Scan(&lb.ID, &lb.Name, &lb.lat, &lb.lon, &lb.alt, &lb.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		launchBases = append(launchBases, &lb)
	}
	return launchBases, nil
}

func (storage *PGstorage) getLaunchBasesQuery(IDs []uint64) squirrel.Sqlizer {
	q := squirrel.Select(LaunchBaseIDColumn, LaunchBaseNameColumn, LaunchBaseLatColumn, LaunchBaseLonColumn, LaunchBaseAltColumn, LaunchBaseCreatedAtColumn).From(launchBaseTableName).
		Where(squirrel.Eq{LaunchBaseIDColumn: IDs}).PlaceholderFormat(squirrel.Dollar)
	return q
}

func (storage *PGstorage) GetDronesByIDs(ctx context.Context, IDs []uint64) ([]*Drone, error) {
	query := storage.getDronesQuery(IDs)
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()
	var drones []*Drone
	for rows.Next() {
		var d Drone
		if err := rows.Scan(&d.ID, &d.serial, &d.model, &d.status, &d.launchbase, &d.CreatedAt); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		drones = append(drones, &d)
	}
	return drones, nil
}

func (storage *PGstorage) getDronesQuery(IDs []uint64) squirrel.Sqlizer {
	q := squirrel.Select(DroneIDColumn, DroneSerialColumn, DroneModelColumn, DroneStatusColumn, DroneLaunchBaseIDColumn, DroneCreatedAtColumn).From(droneTableName).
		Where(squirrel.Eq{DroneIDColumn: IDs}).PlaceholderFormat(squirrel.Dollar)
	return q
}

func (storage *PGstorage) GetMissionsByIDs(ctx context.Context, IDs []uint64) ([]*Mission, error) {
	query := storage.getMissionsQuery(IDs)
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()
	var missions []*Mission
	for rows.Next() {
		var m Mission
		if err := rows.Scan(&m.ID, &m.operator_id, &m.launch_base_id, &m.status, &m.destination_lat, &m.destination_lon, &m.destination_alt, &m.payload_kg, &m.created_at); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		missions = append(missions, &m)
	}
	return missions, nil
}

func (storage *PGstorage) getMissionsQuery(IDs []uint64) squirrel.Sqlizer {
	q := squirrel.Select(MissionIDColumn, MissionOperatorIDColumn, MissionLaunchBaseIDColumn, MissionStatusColumn, MissionDestinationLatColumn, MissionDestinationLonColumn, MissionDestinationAltColumn, MissionPayloadKgColumn, MissionCreatedAtColumn).From(missionTableName).
		Where(squirrel.Eq{MissionIDColumn: IDs}).PlaceholderFormat(squirrel.Dollar)
	return q
}

func (storage *PGstorage) GetMissionDronesByMissionIDs(ctx context.Context, missionIDs []uint64) ([]*MissionDrone, error) {
	query := storage.getMissionDronesQuery(missionIDs)
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()
	var missionDrones []*MissionDrone
	for rows.Next() {
		var md MissionDrone
		if err := rows.Scan(&md.MissionID, &md.DroneID, &md.assigned_by, &md.assigned_at, &md.planned_payload_kg); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		missionDrones = append(missionDrones, &md)
	}
	return missionDrones, nil
}

func (storage *PGstorage) getMissionDronesQuery(missionIDs []uint64) squirrel.Sqlizer {
	q := squirrel.Select(MissionDroneMissionIDColumn, MissionDroneDroneIDColumn, MissionDroneAssignedByColumn, MissionDroneAssignedAtColumn, MissionDronePlannedPayloadKgColumn).From(missionDroneTableName).
		Where(squirrel.Eq{MissionDroneMissionIDColumn: missionIDs}).PlaceholderFormat(squirrel.Dollar)
	return q
}
