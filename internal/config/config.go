package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT value: %v", err)
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "pihole.local"),
		DBPort:     port,
		DBUser:     getEnv("DB_USER", "casaos"),
		DBPassword: getEnv("DB_PASSWORD", "casaos"),
		DBName:     getEnv("DB_NAME", "garage-web"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 