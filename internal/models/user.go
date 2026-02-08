package models

import (
	"time"

	"gorm.io/gorm"
)

// User database model
type User struct {
	ID        uint
	Email     string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Name      string `gorm:"not null"`
	Role      string `gorm:"default:user"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
