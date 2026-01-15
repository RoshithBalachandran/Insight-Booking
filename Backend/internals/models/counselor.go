package models

type Counselor struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"not null"`
	Qualification  string
	Specialization string
	ProfileImage   string
	Email          string `json:"email"`
	Password       string `json:"password"`
	IsActive       bool   `gorm:"default:true"`
}
