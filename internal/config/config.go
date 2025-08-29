package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            int
	Name            string
	Username        string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

var (
	instance *Config
	once     sync.Once
)

// GetConfig returns the singleton configuration instance
func GetConfig() *Config {
	once.Do(func() {
		instance = loadConfig()
	})
	return instance
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg := &Config{
		Database: DatabaseConfig{
			Driver:          getEnv("DATABASE_DRIVER", "postgres"),
			Host:            getEnv("DATABASE_HOST", "localhost"),
			Port:            getEnvAsInt("DATABASE_PORT", 5432),
			Name:            getEnv("DATABASE_NAME", "sample_grpc_server"),
			Username:        getEnv("DATABASE_USERNAME", "postgres"),
			Password:        getEnv("DATABASE_PASSWORD", "postgres"),
			SSLMode:         getEnv("DATABASE_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DATABASE_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: time.Duration(getEnvAsInt("DATABASE_CONN_MAX_LIFETIME", 300)) * time.Second,
		},
	}

	return cfg
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.Name,
		c.SSLMode,
	)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}