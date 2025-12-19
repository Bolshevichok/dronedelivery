package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PGstorage struct {
	db *pgxpool.Pool
}

func NewPGStorge(connString string) (*PGstorage, error) {

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения")
	}
	storage := &PGstorage{
		db: db,
	}
	err = storage.initTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PGstorage) initTables() error {
	queries := []string{
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			%s SERIAL PRIMARY KEY,
			%s VARCHAR(255) UNIQUE NOT NULL,
			%s VARCHAR(100) NOT NULL,
			%s TIMESTAMP NOT NULL DEFAULT NOW()
		)`, operatorTableName, OperatorIDColumn, OperatorEmailCol, OperatorNameCol, OperatorCreatedAtCol),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			%s SERIAL PRIMARY KEY,
			%s VARCHAR(100) NOT NULL,
			%s FLOAT NOT NULL,
			%s FLOAT NOT NULL,
			%s FLOAT NOT NULL,
			%s TIMESTAMP NOT NULL DEFAULT NOW()
		)`, launchBaseTableName, LaunchBaseIDColumn, LaunchBaseNameColumn, LaunchBaseLatColumn, LaunchBaseLonColumn, LaunchBaseAltColumn, LaunchBaseCreatedAtColumn),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			%s SERIAL PRIMARY KEY,
			%s VARCHAR(100) UNIQUE NOT NULL,
			%s VARCHAR(100) NOT NULL,
			%s VARCHAR(50) NOT NULL,
			%s BIGINT REFERENCES %s(%s),
			%s TIMESTAMP NOT NULL DEFAULT NOW()
		)`, droneTableName, DroneIDColumn, DroneSerialColumn, DroneModelColumn, DroneStatusColumn, DroneLaunchBaseIDColumn, launchBaseTableName, LaunchBaseIDColumn, DroneCreatedAtColumn),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			%s SERIAL PRIMARY KEY,
			%s BIGINT REFERENCES %s(%s),
			%s BIGINT REFERENCES %s(%s),
			%s VARCHAR(50) NOT NULL,
			%s FLOAT NOT NULL,
			%s FLOAT NOT NULL,
			%s FLOAT NOT NULL,
			%s FLOAT NOT NULL,
			%s TIMESTAMP NOT NULL DEFAULT NOW()
		)`, missionTableName, MissionIDColumn, MissionOperatorIDColumn, operatorTableName, OperatorIDColumn, MissionLaunchBaseIDColumn, launchBaseTableName, LaunchBaseIDColumn, MissionStatusColumn, MissionDestinationLatColumn, MissionDestinationLonColumn, MissionDestinationAltColumn, MissionPayloadKgColumn, MissionCreatedAtColumn),

		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			%s BIGINT REFERENCES %s(%s),
			%s BIGINT REFERENCES %s(%s),
			%s BIGINT REFERENCES %s(%s),
			%s TIMESTAMP NOT NULL DEFAULT NOW(),
			%s FLOAT NOT NULL,
			PRIMARY KEY (%s, %s)
		)`, missionDroneTableName, MissionDroneMissionIDColumn, missionTableName, MissionIDColumn, MissionDroneDroneIDColumn, droneTableName, DroneIDColumn, MissionDroneAssignedByColumn, operatorTableName, OperatorIDColumn, MissionDroneAssignedAtColumn, MissionDronePlannedPayloadKgColumn, MissionDroneMissionIDColumn, MissionDroneDroneIDColumn),
	}

	for _, sql := range queries {
		_, err := s.db.Exec(context.Background(), sql)
		if err != nil {
			return errors.Wrap(err, "init tables error")
		}
	}
	return nil
}
