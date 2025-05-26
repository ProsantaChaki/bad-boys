package main

import (
	"bad_boyes/internal/config"
	"bad_boyes/internal/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize router
	r := gin.Default()

	// Setup routes
	router.SetupRoutes(r)

	// Start server
	if err := r.Run(":3003"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
