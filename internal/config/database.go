package config

import (
	"fiber-api-boilerplate/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB creates database connection and runs migrations
func ConnectDB(cfg *Config) *gorm.DB {
	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	// Set log level based on environment
	logLevel := logger.Silent
	if cfg.IsDevelopment() {
		logLevel = logger.Info
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models (development only)
	if cfg.IsDevelopment() {
		if err := db.AutoMigrate(&models.User{}); err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
		log.Println("Database connected and migrated successfully")
	} else {
		log.Println("Database connected (production mode - use golang-migrate for migrations)")
	}

	return db
}
