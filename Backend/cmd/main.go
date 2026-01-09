package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/insight/database"
	"github.com/insight/routes"
	"github.com/joho/godotenv"
)

func main() {

	// env load
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("failed to load env")
	}
	//database connection
	database.Connect()
	r := gin.Default()
	//port Setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//setup Routes
	routes.SetupRouter(r)

	r.Run(":" + port)
}
