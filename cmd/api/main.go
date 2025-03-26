package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"task-manager/internal/handlers"
	"task-manager/internal/middlewares"
	"task-manager/internal/models"
	"task-manager/pkg/database"
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

	// Initialize Gin router
	router := gin.New()

	// Apply middlewares
	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddleware())

	// Setup API routes
	api := router.Group("/api")
	{
		// Public routes (no authentication required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Protected routes (authentication required)
		tasks := api.Group("/tasks")
		tasks.Use(middlewares.AuthMiddleware())
		{
			tasks.POST("/", handlers.CreateTask)
			tasks.GET("/", handlers.GetTasks)
			tasks.GET("/:id", handlers.GetTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.PATCH("/:id/status", handlers.UpdateTaskStatus)
			tasks.DELETE("/:id", handlers.DeleteTask)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

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