package models

import "time"

// RefreshToken database model
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"uniqueIndex;type:text;not null"`
	UserID    uint      `gorm:"index;not null"`
	ExpiredAt time.Time `gorm:"index;not null"`
	CreatedAt time.Time
}
