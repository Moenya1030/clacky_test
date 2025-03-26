package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"task-manager/internal/middlewares"
	"task-manager/internal/models"
	"task-manager/internal/routes"
	"task-manager/pkg/database"
	"task-manager/pkg/utils"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	// Set Gin mode based on environment
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup database models and migrations
	if err := models.SetupModels(db); err != nil {
		log.Fatalf("Failed to setup database models: %v", err)
	}

	// Schedule periodic session cleanup
	go scheduleSessionCleanup()

	// Initialize Gin router
	router := gin.New()

	// Apply middlewares
	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddleware())

	// Setup routes using the routes package
	routes.SetupRoutes(router)

	// Get port from environment variable, default to 8080
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// scheduleSessionCleanup runs the session cleanup process periodically
func scheduleSessionCleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running session cleanup")
		utils.CleanupSessions()
	}
}