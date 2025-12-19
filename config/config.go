package config

import (
	"fmt"
	"os"

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
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}
