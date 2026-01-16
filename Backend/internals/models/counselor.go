package models

type Counselor struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"not null"`
	Qualification  string `json:"qualification"`
	Specialization string `json:"specialization"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Block          bool   `json:"block"`
	IsActive       bool   `gorm:"default:true"`
}
