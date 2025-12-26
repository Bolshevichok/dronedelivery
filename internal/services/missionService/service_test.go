package missionService

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/models"
	mockStorage "github.com/Bolshevichok/dronedelivery/internal/services/missionService/mocks"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MissionServiceTestSuite struct {
	suite.Suite
	mockStorage     *mockStorage.MissionStorage
	mockKafkaWriter *kafka.Writer
	mockRedisClient *redis.Client
	service         *MissionService
}

func (suite *MissionServiceTestSuite) SetupTest() {
	suite.mockStorage = &mockStorage.MissionStorage{}
	suite.mockKafkaWriter = &kafka.Writer{}
	suite.mockRedisClient = &redis.Client{}

	cfg := &config.Config{
		Kafka: config.KafkaConfig{
			Host:                 "localhost",
			Port:                 9092,
			MissionsCreatedTopic: "missions.created",
		},
		Redis: config.RedisConfig{
			Host: "localhost",
			Port: 6379,
		},
	}

	suite.service = NewMissionService(context.Background(), suite.mockStorage, cfg)
}

func (suite *MissionServiceTestSuite) TearDownTest() {
	suite.mockStorage.AssertExpectations(suite.T())
}

func (suite *MissionServiceTestSuite) TestCreateMission() {

	expectedMission := &models.Mission{
		ID:             1,
		OperatorID:     1,
		LaunchBaseID:   1,
		Status:         "created",
		DestinationLat: 55.7558,
		DestinationLon: 37.6173,
		DestinationAlt: 100,
		PayloadKg:      1.2,
		CreatedAt:      time.Now(),
	}

	suite.mockStorage.On("UpsertMissions", mock.Anything, mock.Anything).Return([]*models.Mission{expectedMission}, nil)

	missions := []*models.Mission{
		{
			OperatorID:     1,
			LaunchBaseID:   1,
			Status:         "",
			DestinationLat: 55.7558,
			DestinationLon: 37.6173,
			DestinationAlt: 100,
			PayloadKg:      1.2,
			CreatedAt:      time.Now(),
		},
	}

	_, err := suite.service.UpsertMissions(context.Background(), missions)

	assert.NoError(suite.T(), err)
}

func (suite *MissionServiceTestSuite) TestUpsertMissionsError() {
	missions := []*models.Mission{
		{
			OperatorID:     1,
			LaunchBaseID:   1,
			DestinationLat: 55.7558,
			DestinationLon: 37.6173,
			DestinationAlt: 100,
			PayloadKg:      1.2,
		},
	}

	suite.mockStorage.On("UpsertMissions", mock.Anything, mock.Anything).Return([]*models.Mission(nil), fmt.Errorf("storage error"))

	_, err := suite.service.UpsertMissions(context.Background(), missions)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to upsert mission")
}

func (suite *MissionServiceTestSuite) TestUpsertMissionsKafkaPublishErrorIsIgnored() {
	expectedMission := &models.Mission{ID: 1, Status: "created"}
	suite.mockStorage.On("UpsertMissions", mock.Anything, mock.Anything).Return([]*models.Mission{expectedMission}, nil)

	// Force publish failure to hit the slog.Error branch.
	suite.service.kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP("127.0.0.1:1"),
		Topic:    "missions.created",
		Balancer: &kafka.LeastBytes{},
	}

	_, err := suite.service.UpsertMissions(context.Background(), []*models.Mission{{Status: "created"}})
	assert.NoError(suite.T(), err)
}

func (suite *MissionServiceTestSuite) TestGetMission() {
	expectedMission := &models.Mission{
		ID:             1,
		OperatorID:     1,
		LaunchBaseID:   1,
		Status:         "created",
		DestinationLat: 55.7558,
		DestinationLon: 37.6173,
		DestinationAlt: 100,
		PayloadKg:      1.2,
		CreatedAt:      time.Now(),
		Operator: models.Operator{
			ID:        1,
			Email:     "test@example.com",
			Name:      "Test Operator",
			CreatedAt: time.Now(),
		},
		LaunchBase: models.LaunchBase{
			ID:        1,
			Name:      "Test Base",
			Lat:       55.7558,
			Lon:       37.6173,
			Alt:       100,
			CreatedAt: time.Now(),
		},
	}

	suite.mockStorage.On("GetMissionsByIDs", mock.Anything, []uint64{1}).Return([]*models.Mission{expectedMission}, nil)
	suite.mockStorage.On("GetOperatorsByIDs", mock.Anything, []uint64{1}).Return([]*models.Operator{&expectedMission.Operator}, nil)
	suite.mockStorage.On("GetLaunchBasesByIDs", mock.Anything, []uint64{1}).Return([]*models.LaunchBase{&expectedMission.LaunchBase}, nil)
	suite.mockStorage.On("GetMissionDronesByMissionID", mock.Anything, uint64(1)).Return([]*models.MissionDrone{}, nil)
	suite.mockStorage.On("GetDronesByIDs", mock.Anything, []uint64{}).Return([]*models.Drone{}, nil)

	mission, err := suite.service.GetMission(context.Background(), 1)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), mission)
	assert.Equal(suite.T(), uint64(1), mission.ID)
}

func (suite *MissionServiceTestSuite) TestGetMissionWithDrones() {
	expectedMission := &models.Mission{
		ID:           1,
		OperatorID:   1,
		LaunchBaseID: 1,
		Status:       "created",
	}
	operator := &models.Operator{ID: 1}
	launchBase := &models.LaunchBase{ID: 1}
	missionDrones := []*models.MissionDrone{{MissionID: 1, DroneID: 7}}
	drones := []*models.Drone{{ID: 7}}

	suite.mockStorage.On("GetMissionsByIDs", mock.Anything, []uint64{1}).Return([]*models.Mission{expectedMission}, nil)
	suite.mockStorage.On("GetOperatorsByIDs", mock.Anything, []uint64{1}).Return([]*models.Operator{operator}, nil)
	suite.mockStorage.On("GetLaunchBasesByIDs", mock.Anything, []uint64{1}).Return([]*models.LaunchBase{launchBase}, nil)
	suite.mockStorage.On("GetMissionDronesByMissionID", mock.Anything, uint64(1)).Return(missionDrones, nil)
	suite.mockStorage.On("GetDronesByIDs", mock.Anything, []uint64{7}).Return(drones, nil)

	mission, err := suite.service.GetMission(context.Background(), 1)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), mission)
	assert.Len(suite.T(), mission.Drones, 1)
	assert.Equal(suite.T(), uint64(7), mission.Drones[0].ID)
}

func (suite *MissionServiceTestSuite) TestGetMissionOperatorError() {
	suite.mockStorage.On("GetMissionsByIDs", mock.Anything, []uint64{1}).Return([]*models.Mission{{ID: 1, OperatorID: 1, LaunchBaseID: 1}}, nil)
	suite.mockStorage.On("GetOperatorsByIDs", mock.Anything, []uint64{1}).Return([]*models.Operator(nil), fmt.Errorf("storage error"))

	_, err := suite.service.GetMission(context.Background(), 1)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get operator")
}

func (suite *MissionServiceTestSuite) TestGetMissionLaunchBaseError() {
	suite.mockStorage.On("GetMissionsByIDs", mock.Anything, []uint64{1}).Return([]*models.Mission{{ID: 1, OperatorID: 1, LaunchBaseID: 1}}, nil)
	suite.mockStorage.On("GetOperatorsByIDs", mock.Anything, []uint64{1}).Return([]*models.Operator{}, nil)
	suite.mockStorage.On("GetLaunchBasesByIDs", mock.Anything, []uint64{1}).Return([]*models.LaunchBase(nil), fmt.Errorf("storage error"))

	_, err := suite.service.GetMission(context.Background(), 1)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get launch base")
}

func (suite *MissionServiceTestSuite) TestGetMissionMissionDronesError() {
	suite.mockStorage.On("GetMissionsByIDs", mock.Anything, []uint64{1}).Return([]*models.Mission{{ID: 1, OperatorID: 1, LaunchBaseID: 1}}, nil)
	suite.mockStorage.On("GetOperatorsByIDs", mock.Anything, []uint64{1}).Return([]*models.Operator{}, nil)
	suite.mockStorage.On("GetLaunchBasesByIDs", mock.Anything, []uint64{1}).Return([]*models.LaunchBase{}, nil)
	suite.mockStorage.On("GetMissionDronesByMissionID", mock.Anything, uint64(1)).Return([]*models.MissionDrone(nil), fmt.Errorf("storage error"))

	_, err := suite.service.GetMission(context.Background(), 1)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get mission drones")
}

func (suite *MissionServiceTestSuite) TestGetMissionDronesError() {
	suite.mockStorage.On("GetMissionsByIDs", mock.Anything, []uint64{1}).Return([]*models.Mission{{ID: 1, OperatorID: 1, LaunchBaseID: 1}}, nil)
	suite.mockStorage.On("GetOperatorsByIDs", mock.Anything, []uint64{1}).Return([]*models.Operator{}, nil)
	suite.mockStorage.On("GetLaunchBasesByIDs", mock.Anything, []uint64{1}).Return([]*models.LaunchBase{}, nil)
	suite.mockStorage.On("GetMissionDronesByMissionID", mock.Anything, uint64(1)).Return([]*models.MissionDrone{{MissionID: 1, DroneID: 9}}, nil)
	suite.mockStorage.On("GetDronesByIDs", mock.Anything, []uint64{9}).Return([]*models.Drone(nil), fmt.Errorf("storage error"))

	_, err := suite.service.GetMission(context.Background(), 1)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get drones")
}

func (suite *MissionServiceTestSuite) TestProcessMissionLifecycleNilEvent() {
	err := suite.service.ProcessMissionLifecycle(context.Background(), nil)
	assert.NoError(suite.T(), err)
}

func (suite *MissionServiceTestSuite) TestProcessMissionLifecycleUpdatesStatus() {
	suite.mockStorage.On("UpdateMissionStatus", mock.Anything, uint64(1), "delivered").Return(nil)

	err := suite.service.ProcessMissionLifecycle(context.Background(), &models.MissionLifecycleEvent{
		MissionID: 1,
		Status:    "delivered",
		Timestamp: time.Now(),
	})

	assert.NoError(suite.T(), err)
}

func (suite *MissionServiceTestSuite) TestGetDroneTelemetrySuccess() {
	mr, err := miniredis.Run()
	assert.NoError(suite.T(), err)
	suite.T().Cleanup(mr.Close)

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	suite.service.redisClient = client

	err = client.Set(context.Background(), "drone:123", "{\"ok\":true}", 0).Err()
	assert.NoError(suite.T(), err)

	val, err := suite.service.GetDroneTelemetry(context.Background(), "123")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "{\"ok\":true}", val)
}

func (suite *MissionServiceTestSuite) TestGetDroneTelemetryMissingKey() {
	mr, err := miniredis.Run()
	assert.NoError(suite.T(), err)
	suite.T().Cleanup(mr.Close)

	suite.service.redisClient = redis.NewClient(&redis.Options{Addr: mr.Addr()})

	_, err = suite.service.GetDroneTelemetry(context.Background(), "404")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get telemetry")
}

func (suite *MissionServiceTestSuite) TestGetMissionNotFound() {
	suite.mockStorage.On("GetMissionsByIDs", mock.Anything, []uint64{1}).Return([]*models.Mission{}, nil)

	_, err := suite.service.GetMission(context.Background(), 1)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "mission not found")
}

func (suite *MissionServiceTestSuite) TestGetMissionStorageError() {
	suite.mockStorage.On("GetMissionsByIDs", mock.Anything, []uint64{1}).Return([]*models.Mission(nil), fmt.Errorf("storage error"))

	_, err := suite.service.GetMission(context.Background(), 1)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "storage error")
}

func TestMissionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MissionServiceTestSuite))
}
