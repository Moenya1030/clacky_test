package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
	// ErrMaxRetriesReached is returned when the database connection fails after max retries
	ErrMaxRetriesReached = errors.New("max connection retries reached")
)

// DBConfig holds database connection configuration
type DBConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	Charset        string
	ParseTime      bool
	Loc            string
	MaxOpenConns   int
	MaxIdleConns   int
	ConnMaxLifetime time.Duration
	RetryAttempts  int
	RetryDelay     time.Duration
	AllowNativeAuth bool
	UseSocket      bool
}

// LoadDBConfig loads database configuration from environment variables
func LoadDBConfig() DBConfig {
	// Get database connection parameters from environment variables with defaults
	maxOpenConns, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_OPEN_CONNS", "100"))
	maxIdleConns, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_IDLE_CONNS", "10"))
	connMaxLifetime, _ := time.ParseDuration(getEnvOrDefault("DB_CONN_MAX_LIFETIME", "1h"))
	retryAttempts, _ := strconv.Atoi(getEnvOrDefault("DB_RETRY_ATTEMPTS", "3"))
	retryDelay, _ := time.ParseDuration(getEnvOrDefault("DB_RETRY_DELAY", "2s"))
	allowNativeAuth := strings.ToLower(getEnvOrDefault("DB_ALLOW_NATIVE_AUTH", "true")) == "true"
	useSocket := strings.ToLower(getEnvOrDefault("DB_USE_SOCKET", "false")) == "true"

	// Handle parseTime parameter correctly
	parseTime := true
	dbParseTimeValue := strings.ToLower(getEnvOrDefault("DB_PARSE_TIME", "true"))
	if dbParseTimeValue == "false" || dbParseTimeValue == "0" {
		parseTime = false
	}

	return DBConfig{
		Host:           getEnvOrDefault("DB_HOST", "localhost"),
		Port:           getEnvOrDefault("DB_PORT", "3306"),
		User:           getEnvOrDefault("DB_USER", "root"),
		Password:       getEnvOrDefault("DB_PASSWORD", ""),
		Name:           getEnvOrDefault("DB_NAME", "task_manager"),
		Charset:        getEnvOrDefault("DB_CHARSET", "utf8mb4"),
		ParseTime:      parseTime,
		Loc:            getEnvOrDefault("DB_LOC", "Local"),
		MaxOpenConns:   maxOpenConns,
		MaxIdleConns:   maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
		RetryAttempts:  retryAttempts,
		RetryDelay:     retryDelay,
		AllowNativeAuth: allowNativeAuth,
		UseSocket:      useSocket,
	}
}

// BuildDSN builds database connection string based on configuration
func (c DBConfig) BuildDSN() string {
	var dsn string
	
	// Two connection methods: TCP or Socket
	if c.UseSocket {
		// Socket connection (useful for some deployment environments)
		socketPath := getEnvOrDefault("DB_SOCKET_PATH", "/var/run/mysqld/mysqld.sock")
		dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=%s&parseTime=%t&loc=%s",
			c.User, c.Password, socketPath, c.Name, c.Charset, c.ParseTime, c.Loc)
	} else {
		// Standard TCP connection
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
			c.User, c.Password, c.Host, c.Port, c.Name, c.Charset, c.ParseTime, c.Loc)
	}

	// Add additional connection parameters for better stability and compatibility
	params := []string{
		"timeout=30s",
		"readTimeout=30s",
		"writeTimeout=30s",
	}

	// Add parameter for multiple authentication methods support
	if c.AllowNativeAuth {
		params = append(params, "allowNativePasswords=true")
		params = append(params, "allowOldPasswords=true")
	}

	// Append all parameters to DSN
	return dsn + "&" + strings.Join(params, "&")
}

// InitDB initializes the database connection using environment variables
func InitDB() (*gorm.DB, error) {
	config := LoadDBConfig()

	// Set up GORM logger configuration based on environment
	logLevel := logger.Silent
	if strings.ToLower(getEnvOrDefault("APP_ENV", "production")) == "development" {
		logLevel = logger.Info
	} else if strings.ToLower(getEnvOrDefault("LOG_LEVEL", "")) == "debug" {
		logLevel = logger.Info
	}

	// Configure GORM with logger settings
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	// Initialize database connection with retry mechanism
	var (
		db  *gorm.DB
		err error
	)

	ctx, cancel := context.WithTimeout(context.Background(), 
		time.Duration(config.RetryAttempts)*config.RetryDelay)
	defer cancel()

	log.Printf("Connecting to database %s on %s:%s...", config.Name, config.Host, config.Port)
	
	for attempt := 0; attempt < config.RetryAttempts; attempt++ {
		if attempt > 0 {
			log.Printf("Retrying database connection (attempt %d of %d) in %v...", 
				attempt+1, config.RetryAttempts, config.RetryDelay)
			time.Sleep(config.RetryDelay)
			// Exponential backoff
			config.RetryDelay *= 2
		}

		// Try primary connection method
		dsn := config.BuildDSN()
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
		
		if err == nil {
			break
		}
		
		// If primary connection fails and we're on the last attempt, try socket if not already using it
		if attempt == config.RetryAttempts-1 && !config.UseSocket {
			log.Printf("TCP connection failed, trying socket connection as fallback...")
			socketConfig := config
			socketConfig.UseSocket = true
			fallbackDsn := socketConfig.BuildDSN()
			db, err = gorm.Open(mysql.Open(fallbackDsn), gormConfig)
			if err == nil {
				log.Printf("Socket connection successful!")
				break
			}
		}

		log.Printf("Database connection attempt %d failed: %v", attempt+1, err)
		
		// Check context to see if we've exceeded overall timeout
		if ctx.Err() != nil {
			return nil, fmt.Errorf("database connection timed out: %w", ctx.Err())
		}
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMaxRetriesReached, err)
	}

	// Store the global DB instance
	DB = db

	// Configure connection pool settings
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %v", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Print diagnostic information
	if err := printDatabaseInfo(sqlDB); err != nil {
		log.Printf("WARNING: Could not retrieve database information: %v", err)
	}

	log.Printf("Successfully connected to database %s", config.Name)
	return DB, nil
}

// printDatabaseInfo prints diagnostic information about the database
func printDatabaseInfo(db *sql.DB) error {
	var version string
	err := db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		return err
	}
	log.Printf("Connected to MySQL server version: %s", version)
	
	// Get additional database variables if in development mode
	if strings.ToLower(getEnvOrDefault("APP_ENV", "production")) == "development" {
		rows, err := db.Query("SHOW VARIABLES WHERE Variable_name IN " +
			"('max_connections', 'wait_timeout', 'interactive_timeout', 'default_authentication_plugin')")
		if err != nil {
			return err
		}
		defer rows.Close()
		
		log.Println("MySQL configuration:")
		for rows.Next() {
			var name, value string
			if err := rows.Scan(&name, &value); err != nil {
				return err
			}
			log.Printf("  %s: %s", name, value)
		}
	}
	
	// Test ping to verify connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	
	return nil
}

// getEnvOrDefault retrieves an environment variable or returns a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}