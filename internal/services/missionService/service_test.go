package missionService

import (
	"context"
	"testing"
	"time"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/mocks"
	"github.com/Bolshevichok/dronedelivery/internal/models"
	mission_api "github.com/Bolshevichok/dronedelivery/internal/pb/mission_api"
	models_pb "github.com/Bolshevichok/dronedelivery/internal/pb/models"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MissionServiceTestSuite struct {
	suite.Suite
	mockStorage     *mocks.Storage
	mockKafkaWriter *kafka.Writer
	mockRedisClient *redis.Client
	service         MissionService
}

func (suite *MissionServiceTestSuite) SetupTest() {
	suite.mockStorage = &mocks.Storage{}
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

func TestMissionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MissionServiceTestSuite))
}
