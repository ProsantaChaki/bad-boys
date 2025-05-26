package main

import (
	"bad_boyes/internal/routes"
	"bad_boyes/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize router
	r := gin.Default()

	// Setup routes
	routes.SetupAuthRoutes(r)

	// Start server
	if err := r.Run(":3002"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
