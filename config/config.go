package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Logging  LoggingConfig
}

// AppConfig contains application-related configuration
type AppConfig struct {
	Port string
	Env  string
}

// DatabaseConfig contains database-related configuration
type DatabaseConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
	Charset   string
	ParseTime bool
	Loc       string
}

// JWTConfig contains JWT-related configuration
type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

// LoggingConfig contains logging-related configuration
type LoggingConfig struct {
	Level string
}

var config *Config

// Load initializes the configuration
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// Initialize config singleton if not already initialized
	if config == nil {
		config = &Config{
			App: AppConfig{
				Port: getEnvOrDefault("APP_PORT", "8080"),
				Env:  getEnvOrDefault("APP_ENV", "development"),
			},
			Database: DatabaseConfig{
				Host:      getEnvOrDefault("DB_HOST", "localhost"),
				Port:      getEnvOrDefault("DB_PORT", "3306"),
				User:      getEnvOrDefault("DB_USER", "root"),
				Password:  getEnvOrDefault("DB_PASSWORD", ""),
				Name:      getEnvOrDefault("DB_NAME", "task_manager"),
				Charset:   getEnvOrDefault("DB_CHARSET", "utf8mb4"),
				ParseTime: getBoolEnvOrDefault("DB_PARSE_TIME", true),
				Loc:       getEnvOrDefault("DB_LOC", "Local"),
			},
			JWT: JWTConfig{
				Secret:    getEnvOrDefault("JWT_SECRET", "default_jwt_secret_change_me"),
				ExpiresIn: getDurationEnvOrDefault("JWT_EXPIRES_IN", 24*time.Hour),
			},
			Logging: LoggingConfig{
				Level: getEnvOrDefault("LOG_LEVEL", "info"),
			},
		}
	}

	return config
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	if config == nil {
		return Load()
	}
	return config
}

// IsProduction returns true if the application is running in production mode
func IsProduction() bool {
	return GetConfig().App.Env == "production"
}

// IsDevelopment returns true if the application is running in development mode
func IsDevelopment() bool {
	return GetConfig().App.Env == "development"
}

// Helper functions for retrieving environment variables with defaults

// getEnvOrDefault retrieves an environment variable or returns a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getBoolEnvOrDefault retrieves a boolean environment variable or returns a default value if not set
func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Convert string value to boolean
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Warning: Could not parse %s as boolean, using default: %v", key, defaultValue)
		return defaultValue
	}
	return boolValue
}

// getIntEnvOrDefault retrieves an integer environment variable or returns a default value if not set
func getIntEnvOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// Convert string value to integer
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Warning: Could not parse %s as integer, using default: %d", key, defaultValue)
		return defaultValue
	}
	return intValue
}

// getDurationEnvOrDefault retrieves a duration environment variable or returns a default value if not set
func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	// If value doesn't contain a unit (like "h", "m", "s"), assume hours
	if !strings.ContainsAny(value, "hms") {
		value = value + "h"
	}

	// Convert string value to duration
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Warning: Could not parse %s as duration, using default: %v", key, defaultValue)
		return defaultValue
	}
	return duration
}