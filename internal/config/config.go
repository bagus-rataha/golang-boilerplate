package config

import (
	"log"
	"os"
	"strconv"
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
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExpire  time.Duration
	JWTRefreshExpire time.Duration
	Port             string
	AppEnv           string
	AllowedOrigins   string
	RateLimitMax     int
	RateLimitWindow  time.Duration
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

	// Parse rate limit configuration
	rateLimitMax, err := strconv.Atoi(getEnv("RATE_LIMIT_MAX", "100"))
	if err != nil {
		log.Fatal("Invalid RATE_LIMIT_MAX format")
	}

	rateLimitWindow, err := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW", "1m"))
	if err != nil {
		log.Fatal("Invalid RATE_LIMIT_WINDOW format")
	}

	return &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "fiber_api"),
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "change-this-access-secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "change-this-refresh-secret"),
		JWTAccessExpire:  accessExpire,
		JWTRefreshExpire: refreshExpire,
		Port:             getEnv("PORT", "8000"),
		AppEnv:           getEnv("APP_ENV", "development"),
		AllowedOrigins:   getEnv("ALLOWED_ORIGINS", ""),
		RateLimitMax:     rateLimitMax,
		RateLimitWindow:  rateLimitWindow,
	}
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
