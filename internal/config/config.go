package config

import (
	"fmt"
	"os"
)

type Config struct {
	WeatherAPIKey string
	Database      DatabaseConfig
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	DSN      string
}

func Load() *Config {
	cfg := &Config{
		WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		Database: DatabaseConfig{
			User:     getEnvOrDefault("DB_USER", "root"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "3306"),
			Name:     getEnvOrDefault("DB_NAME", "weather_db"),
		},
	}

	cfg.Database.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
