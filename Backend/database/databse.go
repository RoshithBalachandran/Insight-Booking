package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/insight/internals/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable Timezone=Asia/Kolkata",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := DB.AutoMigrate(&models.User{},
		&models.RefreshToken{},
		&models.Service{},
		&models.Counselor{},
	); err != nil {
		log.Fatal("failed to automigrate:", err)
	}

	log.Println("Database connected successfully")
}
