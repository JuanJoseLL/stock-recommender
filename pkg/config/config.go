package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	API      APIConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type APIConfig struct {
	URL string
	Key string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 26257),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "stockmarket"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		API: APIConfig{
			URL: getEnv("API_URL", "https://8j5baasof2.execute-api.us-west-2.amazonaws.com"),
			Key: getEnv("API_KEY", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdHRlbXB0cyI6MSwiZW1haWwiOiJqdWFuam9sbzEyQGhvdG1haWwuY29tIiwiZXhwIjoxNzUwMTg2MTk2LCJpZCI6IjAiLCJwYXNzd29yZCI6Iicgb3IgMT0xIG9yICcnPScifQ.i3iwFvk5BA5l_AonHDYb4RhEJSZfha3nVsb4ZFvQ7u0"),
		},
	}
}

func (c *Config) GetDatabaseURL() string {
	if c.Database.Password == "" {
		return fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=%s",
			c.Database.User,
			c.Database.Host,
			c.Database.Port,
			c.Database.DBName,
			c.Database.SSLMode,
		)
	}
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
