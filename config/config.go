package config

import (
	"fmt"
	"os"
	"strconv"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	Database       DatabaseConfig       `yaml:"database"`
	Kafka          KafkaConfig          `yaml:"kafka"`
	Redis          RedisConfig          `yaml:"redis"`
	MissionService MissionServiceConfig `yaml:"mission_service"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type KafkaConfig struct {
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	MissionsCreatedTopic   string `yaml:"missions_created_topic"`
	MissionsLifecycleTopic string `yaml:"missions_lifecycle_topic"`
	DroneLifecycleTopic    string `yaml:"drone_lifecycle_topic"`
	DroneTelemetryTopic    string `yaml:"drone_telemetry_topic"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MissionServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func LoadConfig(filename string) (*Config, error) {
	if filename == "" || filename == "env" {
		return loadConfigFromEnv()
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		// Fallback to env if file not found
		return loadConfigFromEnv()
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}

func loadConfigFromEnv() (*Config, error) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnvInt("DATABASE_PORT", 5432),
			Username: getEnv("DATABASE_USERNAME", "admin"),
			Password: getEnv("DATABASE_PASSWORD", "admin"),
			DBName:   getEnv("DATABASE_NAME", "dronedelivery"),
			SSLMode:  getEnv("DATABASE_SSL_MODE", "disable"),
		},
		Kafka: KafkaConfig{
			Host:                   getEnv("KAFKA_HOST", "localhost"),
			Port:                   getEnvInt("KAFKA_PORT", 9092),
			MissionsCreatedTopic:   getEnv("KAFKA_MISSIONS_CREATED_TOPIC", "missions.created"),
			MissionsLifecycleTopic: getEnv("KAFKA_MISSIONS_LIFECYCLE_TOPIC", "missions.lifecycle"),
			DroneLifecycleTopic:    getEnv("KAFKA_DRONE_LIFECYCLE_TOPIC", "drone.lifecycle"),
			DroneTelemetryTopic:    getEnv("KAFKA_DRONE_TELEMETRY_TOPIC", "drone.telemetry"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnvInt("REDIS_PORT", 6379),
		},
		MissionService: MissionServiceConfig{
			Host: getEnv("MISSION_SERVICE_HOST", "localhost"),
			Port: getEnvInt("MISSION_SERVICE_PORT", 8080),
		},
	}
	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
