package routes

import (
	"github.com/gin-gonic/gin"

	"task-manager/internal/handlers"
	"task-manager/internal/middlewares"
)

// SetupRoutes configures all the API routes for the application
func SetupRoutes(router *gin.Engine) {
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
}