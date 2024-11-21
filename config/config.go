package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type DB struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
}

type Config struct {
	Dsn        string `yaml:"dsn"`
	Redis_Addr string `yaml:"redis_addr"`
	Redis_DB   int    `yaml:"redis_db"`
	Port       string `yaml:"port"`
	DB         DB
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &config, nil
}
