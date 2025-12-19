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
		q = q.Values(lb.Name, lb.Lat, lb.Lon, lb.Alt, lb.CreatedAt)
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
		q = q.Values(d.Serial, d.Model, d.Status, d.LaunchBaseID, d.CreatedAt)
	}
	return q
}

func (storage *PGstorage) UpsertMissions(ctx context.Context, missions []*Mission) ([]*Mission, error) {
	query := storage.upsertMissionsQuery(missions)
	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}
	rows, err := storage.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "exec query error")
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var id uint64
		err := rows.Scan(&id)
		if err != nil {
			return nil, errors.Wrap(err, "scan id error")
		}
		missions[i].ID = id
		i++
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows error")
	}
	return missions, nil
}

func (storage *PGstorage) upsertMissionsQuery(missions []*Mission) squirrel.Sqlizer {
	q := squirrel.Insert(missionTableName).Columns(MissionOperatorIDColumn, MissionLaunchBaseIDColumn, MissionStatusColumn, MissionDestinationLatColumn, MissionDestinationLonColumn, MissionDestinationAltColumn, MissionPayloadKgColumn, MissionCreatedAtColumn).
		PlaceholderFormat(squirrel.Dollar)
	for _, m := range missions {
		q = q.Values(m.OperatorID, m.LaunchBaseID, m.Status, m.DestinationLat, m.DestinationLon, m.DestinationAlt, m.PayloadKg, m.CreatedAt)
	}
	q = q.Suffix("RETURNING id")
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
		q = q.Values(md.MissionID, md.DroneID, md.AssignedBy, md.AssignedAt, md.PlannedPayloadKg)
	}
	q = q.Suffix("ON CONFLICT (mission_id, drone_id) DO NOTHING")
	return q
}

func (storage *PGstorage) UpdateMissionStatus(ctx context.Context, missionID uint64, status string) error {
	query := squirrel.Update(missionTableName).Set(MissionStatusColumn, status).Where(squirrel.Eq{MissionIDColumn: missionID}).PlaceholderFormat(squirrel.Dollar)
	queryText, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "generate query error")
	}
	_, err = storage.db.Exec(ctx, queryText, args...)
	if err != nil {
		return errors.Wrap(err, "exec query error")
	}
	return nil
}
