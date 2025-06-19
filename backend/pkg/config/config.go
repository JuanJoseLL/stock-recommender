package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	API         APIConfig
	AlphaVantage AlphaVantageConfig
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

type AlphaVantageConfig struct {
	APIKey string
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
			URL: getEnv("API_URL", "https://localhost:8080"),
			Key: getEnv("API_KEY", ""),
		},
		AlphaVantage: AlphaVantageConfig{
			APIKey: getEnv("ALPHA_VANTAGE_API", ""),
		},
	}
}

func (c *Config) GetDatabaseURL() string {
	if dbURL := getEnv("DATABASE_URL", ""); dbURL != "" {
		fmt.Println("Using DATABASE_URL from environment")
		return dbURL
	}
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
