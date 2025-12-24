package missionService

import (
	"context"
	"testing"
	"time"

	"github.com/Bolshevichok/dronedelivery/config"
	"github.com/Bolshevichok/dronedelivery/internal/mocks"
	"github.com/Bolshevichok/dronedelivery/internal/models"
	missionv1 "github.com/Bolshevichok/dronedelivery/internal/pb/mission/v1"
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
	req := &missionv1.CreateMissionRequest{
		OperatorId:     1,
		LaunchBaseId:   1,
		DestinationLat: 55.7558,
		DestinationLon: 37.6173,
		DestinationAlt: 100,
		PayloadKg:      1.2,
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

	suite.mockStorage.On("UpsertMissions", mock.Anything, mock.MatchedBy(func(missions []*models.Mission) bool {
		return len(missions) == 1 && missions[0].OperatorID == 1 && missions[0].LaunchBaseID == 1 &&
			missions[0].Status == "created" && missions[0].DestinationLat == 55.7558 &&
			missions[0].DestinationLon == 37.6173 && missions[0].DestinationAlt == 100 && missions[0].PayloadKg == 1.2
	})).Return([]*models.Mission{expectedMission}, nil)

	resp, err := suite.service.CreateMission(context.Background(), req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), uint64(1), resp.MissionId)
}

func TestMissionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MissionServiceTestSuite))
}
