package models

import "time"

type RefreshToken struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	TokenHash string `gorm:"uniqueIndex"`
	Revoked   bool   `gorm:"default:false"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
