package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	JWTAccessExpire  time.Duration
	JWTRefreshExpire time.Duration
	Port             string
	AppEnv           string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Parse JWT expiration durations
	accessExpire, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRE", "24h"))
	if err != nil {
		log.Fatal("Invalid JWT_ACCESS_EXPIRE format")
	}

	refreshExpire, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRE", "168h"))
	if err != nil {
		log.Fatal("Invalid JWT_REFRESH_EXPIRE format")
	}

	return &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "fiber_api"),
		JWTSecret:        getEnv("JWT_SECRET", "change-this-secret"),
		JWTAccessExpire:  accessExpire,
		JWTRefreshExpire: refreshExpire,
		Port:             getEnv("PORT", "8000"),
		AppEnv:           getEnv("APP_ENV", "development"),
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
