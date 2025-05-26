package config

import (
	"fmt"
	"os"

	"github.com/JunBumHan/copilot-agent-test/internal/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds the application configuration
type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	ServerPort       string
}

// NewConfig creates a new configuration with environment variables or defaults
func NewConfig() *Config {
	return &Config{
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:       getEnv("POSTGRES_DB", "userdb"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
	}
}

// DatabaseURL returns the PostgreSQL connection string
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s ****** dbname=%s sslmode=disable",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB)
}

// ConnectDatabase connects to the database and initializes it
func (c *Config) ConnectDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(c.DatabaseURL()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&store.UserModel{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// Helper function to get an environment variable or a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
