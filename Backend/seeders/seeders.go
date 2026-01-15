package seeders

import (
	"log"
	"os"

	"github.com/insight/database"
	"github.com/insight/internals/constant"
	"github.com/insight/internals/models"
	"github.com/insight/internals/utils"
)

func SeedAdminDetails() {
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || password == "" {
		log.Fatal("ADMIN_EMAIL and ADMIN_PASSWORD must be set")
	}

	// Check if admin already exists
	var count int64
	database.DB.Model(&models.User{}).
		Where("role = ?", constant.Admin).
		Count(&count)

	if count > 0 {
		log.Println("Admin already exists. Skipping seeding.")
		return
	}

	// Hash password
	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal("Password hashing failed:", err)
	}

	admin := models.User{
		FirstName: "Super",
		LastName:  "Admin",
		Email:     email,
		Password:  hash,
		Role:      constant.Admin,
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		log.Fatal("Admin s	eeding failed:", err)
	}

	log.Println("✅ Admin seeded successfully")
}
