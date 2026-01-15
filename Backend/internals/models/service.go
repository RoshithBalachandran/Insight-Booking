package models

type Service struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	IsActive    bool   `gorm:"default:true"`
}
