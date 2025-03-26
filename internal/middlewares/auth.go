package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"task-manager/pkg/utils"
	"task-manager/pkg/database"
	"task-manager/internal/models"
)

// AuthMiddleware authenticates the user by validating JWT token from request header
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

		// Extract token from the Authorization header
		// Format should be "Bearer {token}"
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must be in format: Bearer {token}",
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token cannot be empty",
			})
			c.Abort()
			return
		}

		// Validate the JWT token
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			status := http.StatusUnauthorized
			errorMsg := "Invalid token"
			
			// Provide more specific error messages based on error type
			if errors.Is(err, jwt.ErrTokenExpired) || strings.Contains(err.Error(), "token expired") {
				errorMsg = "Token has expired"
			} else if strings.Contains(err.Error(), "signature") {
				errorMsg = "Invalid token signature"
			} else if strings.Contains(err.Error(), "parsing") {
				errorMsg = "Token format is invalid"
			}
			
			c.JSON(status, gin.H{
				"error": errorMsg,
			})
			c.Abort()
			return
		}

		// Check if user exists in database
		var user models.User
		result := database.GetDB().First(&user, userID)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not found or invalid token",
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