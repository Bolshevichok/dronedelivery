package pgstorage

import (
	"context"

	"github.com/Bolshevichok/dronedelivery/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (storage *PGstorage) UpsertOperators(ctx context.Context, operators []*models.Operator) error {
	query := squirrel.Insert(operatorTableName).Columns(OperatorEmailCol, OperatorNameCol, OperatorCreatedAtCol).
		PlaceholderFormat(squirrel.Dollar)
	for _, op := range operators {
		query = query.Values(op.Email, op.Name, op.CreatedAt)
	}
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

func (storage *PGstorage) UpsertLaunchBases(ctx context.Context, launchBases []*models.LaunchBase) error {
	query := squirrel.Insert(launchBaseTableName).Columns(LaunchBaseNameColumn, LaunchBaseLatColumn, LaunchBaseLonColumn, LaunchBaseAltColumn, LaunchBaseCreatedAtColumn).
		PlaceholderFormat(squirrel.Dollar)
	for _, lb := range launchBases {
		query = query.Values(lb.Name, lb.Lat, lb.Lon, lb.Alt, lb.CreatedAt)
	}
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

func (storage *PGstorage) UpsertDrones(ctx context.Context, drones []*models.Drone) error {
	query := squirrel.Insert(droneTableName).Columns(DroneSerialColumn, DroneModelColumn, DroneStatusColumn, DroneLaunchBaseIDColumn, DroneCreatedAtColumn).
		PlaceholderFormat(squirrel.Dollar)
	for _, d := range drones {
		query = query.Values(d.Serial, d.Model, d.Status, d.LaunchBaseID, d.CreatedAt)
	}
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

func (storage *PGstorage) UpsertMissions(ctx context.Context, missions []*models.Mission) ([]*models.Mission, error) {
	query := squirrel.Insert(missionTableName).Columns(MissionOperatorIDColumn, MissionLaunchBaseIDColumn, MissionStatusColumn, MissionDestinationLatColumn, MissionDestinationLonColumn, MissionDestinationAltColumn, MissionPayloadKgColumn, MissionCreatedAtColumn).
		PlaceholderFormat(squirrel.Dollar)
	for _, m := range missions {
		query = query.Values(m.OperatorID, m.LaunchBaseID, m.Status, m.DestinationLat, m.DestinationLon, m.DestinationAlt, m.PayloadKg, m.CreatedAt)
	}
	query = query.Suffix("RETURNING id")

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

func (storage *PGstorage) UpsertMissionDrones(ctx context.Context, missionDrones []*models.MissionDrone) error {
	query := squirrel.Insert(missionDroneTableName).Columns(MissionDroneMissionIDColumn, MissionDroneDroneIDColumn, MissionDroneAssignedByColumn, MissionDroneAssignedAtColumn, MissionDronePlannedPayloadKgColumn).
		PlaceholderFormat(squirrel.Dollar)
	for _, md := range missionDrones {
		query = query.Values(md.MissionID, md.DroneID, md.AssignedBy, md.AssignedAt, md.PlannedPayloadKg)
	}
	query = query.Suffix("ON CONFLICT (mission_id, drone_id) DO NOTHING")

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

func (storage *PGstorage) CreateOperator(ctx context.Context, operator *models.Operator) (uint64, error) {
	query := squirrel.Insert(operatorTableName).Columns(OperatorEmailCol, OperatorNameCol, OperatorCreatedAtCol).
		Values(operator.Email, operator.Name, operator.CreatedAt).
		Suffix("RETURNING " + OperatorIDColumn).
		PlaceholderFormat(squirrel.Dollar)
	queryText, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "generate query error")
	}
	var id uint64
	err = storage.db.QueryRow(ctx, queryText, args...).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "exec query error")
	}
	return id, nil
}

func (storage *PGstorage) CreateLaunchBase(ctx context.Context, launchBase *models.LaunchBase) (uint64, error) {
	query := squirrel.Insert(launchBaseTableName).Columns(LaunchBaseNameColumn, LaunchBaseLatColumn, LaunchBaseLonColumn, LaunchBaseAltColumn, LaunchBaseCreatedAtColumn).
		Values(launchBase.Name, launchBase.Lat, launchBase.Lon, launchBase.Alt, launchBase.CreatedAt).
		Suffix("RETURNING " + LaunchBaseIDColumn).
		PlaceholderFormat(squirrel.Dollar)
	queryText, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "generate query error")
	}
	var id uint64
	err = storage.db.QueryRow(ctx, queryText, args...).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "exec query error")
	}
	return id, nil
}

func (storage *PGstorage) CreateDrone(ctx context.Context, drone *models.Drone) (uint64, error) {
	query := squirrel.Insert(droneTableName).Columns(DroneSerialColumn, DroneModelColumn, DroneStatusColumn, DroneLaunchBaseIDColumn, DroneCreatedAtColumn).
		Values(drone.Serial, drone.Model, drone.Status, drone.LaunchBaseID, drone.CreatedAt).
		Suffix("RETURNING " + DroneIDColumn).
		PlaceholderFormat(squirrel.Dollar)
	queryText, args, err := query.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "generate query error")
	}
	var id uint64
	err = storage.db.QueryRow(ctx, queryText, args...).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "exec query error")
	}
	return id, nil
}
