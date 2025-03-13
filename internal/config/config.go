package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort     string
	DatabaseURL    string
	BaseURL        string
	ShortCodeLen   int
	DefaultExpiry  time.Duration
	LoggingEnabled bool
}

// New returns a new Config struct
func New() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	shortCodeLen, _ := strconv.Atoi(getEnv("SHORT_CODE_LEN", "6"))
	if shortCodeLen > 10 {
		shortCodeLen = 10 // Ensure it doesn't exceed database column size
	}
	expiryDays, _ := strconv.Atoi(getEnv("DEFAULT_EXPIRY_DAYS", "30"))
	defaultExpiry := time.Duration(expiryDays) * 24 * time.Hour
	loggingEnabled, _ := strconv.ParseBool(getEnv("LOGGING_ENABLED", "true"))

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable"),
		BaseURL:        getEnv("BASE_URL", "http://localhost:8080"),
		ShortCodeLen:   shortCodeLen,
		DefaultExpiry:  defaultExpiry,
		LoggingEnabled: loggingEnabled,
	}
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}