package middlewares

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// LogLevel represents the logging level
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)

// LoggerMiddleware logs HTTP requests with details like method, path, status code, latency time
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get start time
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate response time
		latency := time.Since(startTime)

		// Get request details
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		// Get configured log level from environment
		configuredLevel := getLogLevel()

		// Create log message
		logMsg := fmt.Sprintf("[%s] %s | %3d | %13v | %15s | %s | %s",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			statusCode,
			latency,
			clientIP,
			path,
			userAgent,
		)

		// Log based on status code and configured level
		switch {
		case statusCode >= 500:
			// Server errors are always logged
			errorLog(logMsg)
		case statusCode >= 400:
			// Client errors are logged at warn level and above
			if configuredLevel != ErrorLevel {
				warnLog(logMsg)
			}
		case statusCode >= 300:
			// Redirection responses are logged at info level and above
			if configuredLevel != ErrorLevel && configuredLevel != WarnLevel {
				infoLog(logMsg)
			}
		default:
			// Successful responses are logged at debug level
			if configuredLevel == DebugLevel {
				debugLog(logMsg)
			}
		}
	}
}

// Get log level from environment variable
func getLogLevel() LogLevel {
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		// Default to info if not specified
		return InfoLevel
	}
}

// Log functions for different levels
func debugLog(msg string) {
	fmt.Printf("\033[37m[DEBUG] %s\033[0m\n", msg) // Light gray
}

func infoLog(msg string) {
	fmt.Printf("\033[32m[INFO] %s\033[0m\n", msg) // Green
}

func warnLog(msg string) {
	fmt.Printf("\033[33m[WARN] %s\033[0m\n", msg) // Yellow
}

func errorLog(msg string) {
	fmt.Printf("\033[31m[ERROR] %s\033[0m\n", msg) // Red
}