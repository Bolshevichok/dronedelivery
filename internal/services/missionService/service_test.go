package missionService

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
	models_pb "github.com/Bolshevichok/dronedelivery/internal/pb/models"
	mockStorage "github.com/Bolshevichok/dronedelivery/internal/services/missionService/mocks"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MissionServiceTestSuite struct {
	suite.Suite
	mockStorage     *mockStorage.Storage
	mockKafkaWriter *kafka.Writer
	mockRedisClient *redis.Client
	service         MissionService
}

func (suite *MissionServiceTestSuite) SetupTest() {
	suite.mockStorage = &mockStorage.Storage{}
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
	req := &mission_api.UpsertMissionsRequest{
		Missions: []*models_pb.Mission{
			{
				OpId:    1,
				BaseId:  1,
				Lat:     55.7558,
				Lon:     37.6173,
				Alt:     100,
				Payload: 1.2,
			},
		},
	}

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

	_, err := suite.service.UpsertMissions(context.Background(), req)

	assert.NoError(suite.T(), err)
}

func (suite *MissionServiceTestSuite) TestUpsertMissionsError() {
	req := &mission_api.UpsertMissionsRequest{
		Missions: []*models_pb.Mission{
			{
				OpId:    1,
				BaseId:  1,
				Lat:     55.7558,
				Lon:     37.6173,
				Alt:     100,
				Payload: 1.2,
			},
		},
	}

	suite.mockStorage.On("UpsertMissions", mock.Anything, mock.Anything).Return([]*models.Mission(nil), fmt.Errorf("storage error"))

	_, err := suite.service.UpsertMissions(context.Background(), req)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to upsert mission")
}

func (suite *MissionServiceTestSuite) TestGetMission() {
	req := &mission_api.GetMissionRequest{
		MissionId: 1,
	}

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

	resp, err := suite.service.GetMission(context.Background(), req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp.Mission)
	assert.Equal(suite.T(), uint64(1), resp.Mission.Id)
}

func (suite *MissionServiceTestSuite) TestGetMissionNotFound() {
	req := &mission_api.GetMissionRequest{
		MissionId: 1,
	}

	suite.mockStorage.On("GetMissionsByIDs", mock.Anything, []uint64{1}).Return([]*models.Mission{}, nil)

	_, err := suite.service.GetMission(context.Background(), req)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "mission not found")
}

func (suite *MissionServiceTestSuite) TestGetMissionStorageError() {
	req := &mission_api.GetMissionRequest{
		MissionId: 1,
	}

	suite.mockStorage.On("GetMissionsByIDs", mock.Anything, []uint64{1}).Return([]*models.Mission(nil), fmt.Errorf("storage error"))

	_, err := suite.service.GetMission(context.Background(), req)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to get mission")
}

func TestMissionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MissionServiceTestSuite))
}
