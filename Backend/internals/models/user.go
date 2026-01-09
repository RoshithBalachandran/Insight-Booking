package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName   string `gorm:"size:100"`
	LastName    string `gorm:"size:100"`
	PhoneNumber string `gorm:"size:15;unique"`
	Email       string `gorm:"size:100;uniqueIndex"`
	Password    string `gorm:"not null"`
	Role        string `gorm:"size:20;default:'patient'"`
	IsActive    bool   `gorm:"default:true"`
}
type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"Role"`
}

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}