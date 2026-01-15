package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/internals/redis"
	"github.com/insight/routes"
	"github.com/insight/seeders"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Failed to load env")
	}

	// Connect DB
	database.Connect()

	// Connect Redis
	redis.ConnectRedis()

	//Seed Admin
	seeders.SeedAdminDetails()
	// Setup router
	r := gin.Default()
	//users routes
	routes.SetupRouter(r)
	//Admin Routes
	routes.SetUpAdminRoutes(r)
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
