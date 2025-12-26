package missionService

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/models"
	mockStorage "github.com/Bolshevichok/dronedelivery/internal/services/missionService/mocks"
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

	mission, err := suite.service.GetMission(context.Background(), 1)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), mission)
	assert.Equal(suite.T(), uint64(1), mission.ID)
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
