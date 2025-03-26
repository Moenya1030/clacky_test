package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"task-manager/pkg/utils"
	"task-manager/pkg/database"
	"task-manager/internal/models"
)

// AuthMiddleware authenticates the user by validating session ID from request header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")

		// Check if Authorization header exists
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract session ID directly from the header
		// Remove requirement for "Bearer " prefix
		sessionID := strings.TrimSpace(authHeader)

		// Validate the session ID
		userID, err := utils.ValidateToken(sessionID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid session: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Check if user exists in database
		var user models.User
		result := database.GetDB().First(&user, userID)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found or invalid session",
			})
			c.Abort()
			return
		}

		// Set user ID in context for later use
		c.Set("userID", userID)
		c.Set("user", user)

		// Continue to the next handler
		c.Next()
	}
}

// GetUserID retrieves the current user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetUser retrieves the current user from context
func GetUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}