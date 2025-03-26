# Task Manager API

A comprehensive task management system backend built with Go, providing features for user authentication, task creation, assignment, and tracking.

## Project Overview

The Task Manager API allows users to:
- Register and log in with JWT authentication
- Create, read, update, and delete tasks
- Set task priorities and deadlines
- Update task statuses (todo, in progress, completed)
- Filter and sort tasks based on various criteria

This RESTful API provides a solid backend foundation for task management applications, with clean architecture and performance in mind.

## Technology Stack

- **Language:** Go 1.21.0
- **Web Framework:** Gin
- **Database:** MySQL
- **ORM:** GORM
- **Authentication:** JWT (JSON Web Tokens)
- **Configuration:** Environment variables with godotenv
- **Password Hashing:** bcrypt

## Project Structure

```
task-manager/
├── cmd/
│   └── api/           # Application entrypoints
│       └── main.go    # Main application file
├── config/            # Configuration management
│   └── config.go
├── docs/              # Documentation
│   └── api.md         # API documentation
├── internal/
│   ├── handlers/      # HTTP request handlers
│   │   ├── auth_handler.go
│   │   └── task_handler.go
│   ├── middlewares/   # HTTP middlewares
│   │   ├── auth.go
│   │   └── logger.go
│   ├── models/        # Database models
│   │   ├── setup.go
│   │   ├── task.go
│   │   └── user.go
│   └── services/      # Business logic
│       ├── task_service.go
│       └── user_service.go
├── pkg/
│   ├── database/      # Database connection management
│   │   └── database.go
│   └── utils/         # Utility functions
│       └── jwt.go
├── .env               # Environment variables
├── go.mod             # Go module definition
└── README.md          # Project documentation
```

## Installation and Setup

### Prerequisites

- Go 1.21.0 or later
- MySQL database
- Git

### Steps to Run

1. **Clone the repository**
   ```
   git clone https://github.com/yourusername/task-manager.git
   cd task-manager
   ```

2. **Set up environment variables**
   - Copy the `.env.example` file (if it exists) to `.env`
   - Modify the values in `.env` according to your environment

3. **Install dependencies**
   ```
   go mod download
   ```

4. **Set up the database**
   - Create a MySQL database named `task_manager` (or as specified in your .env file)
   - The application will automatically create the necessary tables on startup

5. **Run the application**
   ```
   go run cmd/api/main.go
   ```

6. **Verify installation**
   - Access the health check endpoint at `http://localhost:8080/health`
   - You should receive a JSON response: `{"status":"ok"}`

## Environment Variables

The application can be configured using the following environment variables in the `.env` file:

### Application Settings
- `APP_PORT`: The port on which the server will run (default: 8080)
- `APP_ENV`: Application environment (development, production)

### Database Settings
- `DB_HOST`: Database host address (default: localhost)
- `DB_PORT`: Database port (default: 3306)
- `DB_USER`: Database username (default: root)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name (default: task_manager)
- `DB_CHARSET`: Database charset (default: utf8mb4)
- `DB_PARSE_TIME`: Parse time values from database (default: true)
- `DB_LOC`: Database timezone (default: Local)

### JWT Settings
- `JWT_SECRET`: Secret key for signing JWT tokens
- `JWT_EXPIRES_IN`: Token expiration time (default: 24h)

### Logging Settings
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

## API Documentation

For detailed API documentation including endpoints, request/response formats, and authentication details, please refer to the [API Documentation](docs/api.md).

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -m 'Add some feature'`
4. Push to the branch: `git push origin feature-name`
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.