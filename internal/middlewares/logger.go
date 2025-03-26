package middlewares

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"task-manager/config"
)

// LogLevel represents the logging level
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
)

// RequestLogData represents structured log data for HTTP requests
type RequestLogData struct {
	Timestamp   string        `json:"timestamp"`
	Method      string        `json:"method"`
	Path        string        `json:"path"`
	StatusCode  int           `json:"status_code"`
	Latency     time.Duration `json:"latency"`
	LatencyHuman string        `json:"latency_human"`
	ClientIP    string        `json:"client_ip"`
	UserAgent   string        `json:"user_agent"`
	RequestID   string        `json:"request_id,omitempty"`
	Error       string        `json:"error,omitempty"`
	QueryParams string        `json:"query_params,omitempty"`
	ReqSize     int           `json:"request_size,omitempty"`
	RespSize    int           `json:"response_size"`
}

// LoggerMiddleware logs HTTP requests with enhanced details
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get start time
		startTime := time.Now()
		
		// Generate or get request ID (if exists)
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
			c.Header("X-Request-ID", requestID)
		}

		// Process request
		c.Next()

		// Calculate response time
		latency := time.Since(startTime)

		// Get request details
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()
		
		// Get request and response size
		reqSize := c.Request.ContentLength
		respSize := c.Writer.Size()

		// Prepare structured log data
		logData := RequestLogData{
			Timestamp:   time.Now().Format(time.RFC3339),
			Method:      method,
			Path:        path,
			StatusCode:  statusCode,
			Latency:     latency,
			LatencyHuman: latency.String(),
			ClientIP:    clientIP,
			UserAgent:   userAgent,
			RequestID:   requestID,
			RespSize:    respSize,
		}

		// Include request size if available
		if reqSize > 0 {
			logData.ReqSize = int(reqSize)
		}
		
		// Include query parameters if present
		if query != "" {
			logData.QueryParams = query
		}

		// Include error if request failed
		if len(c.Errors) > 0 {
			logData.Error = c.Errors.String()
		}

		// Get log level from config (falls back to environment variable)
		configuredLevel := getConfiguredLogLevel()

		// Convert to JSON for structured logging
		logJSON, err := json.Marshal(logData)
		if err != nil {
			// Fallback to simple format if JSON marshal fails
			simpleLog(statusCode, configuredLevel, fmt.Sprintf("%s %s %d %v %s", 
				method, path, statusCode, latency, clientIP))
			return
		}

		// Log based on status code and configured level
		logMessage := string(logJSON)
		
		switch {
		case statusCode >= 500:
			// Server errors are always logged regardless of level
			errorLog(logMessage)
		case statusCode >= 400:
			// Client errors are logged at warn level and above
			if configuredLevel != ErrorLevel {
				warnLog(logMessage)
			}
		case statusCode >= 300:
			// Redirections are logged at info level and above
			if configuredLevel != ErrorLevel && configuredLevel != WarnLevel {
				infoLog(logMessage)
			}
		default:
			// Successful responses are logged based on level
			if configuredLevel == DebugLevel {
				debugLog(logMessage)
			} else if configuredLevel == InfoLevel {
				infoLog(logMessage)
			}
		}
	}
}

// Get configured log level from config or environment variable
func getConfiguredLogLevel() LogLevel {
	// First try to get from config (which also checks environment)
	configLevel := config.GetConfig().Logging.Level
	
	switch configLevel {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		// If config doesn't have a valid value, check environment directly
		return getLogLevelFromEnv()
	}
}

// Get log level directly from environment variable
func getLogLevelFromEnv() LogLevel {
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

// Fallback simple logging when JSON fails
func simpleLog(statusCode int, level LogLevel, msg string) {
	switch {
	case statusCode >= 500:
		errorLog(msg)
	case statusCode >= 400:
		warnLog(msg)
	default:
		if level == DebugLevel {
			debugLog(msg)
		} else {
			infoLog(msg)
		}
	}
}

// Log functions for different levels with color coding
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