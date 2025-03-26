# Task Manager API Documentation

## Overview

This document outlines the RESTful API endpoints for the Task Manager system. The API provides functionality for user authentication and task management with features like filtering, sorting, and pagination.

## Base URL

All API endpoints are prefixed with `/api`. For local development, the base URL is:

```
http://localhost:8080/api
```

## Authentication

The API uses JWT (JSON Web Token) authentication. After logging in or registering, you will receive a token that must be included in all subsequent requests that require authentication.

### How to Authenticate

Include the JWT token in the Authorization header of your requests using the Bearer scheme:

```
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### Authentication

#### Register a New User

- **URL**: `/auth/register`
- **Method**: `POST`
- **Authentication Required**: No
- **Request Body**:
  ```json
  {
    "username": "johndoe",
    "email": "john.doe@example.com",
    "password": "securepassword123"
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john.doe@example.com",
      "created_at": "2023-01-15T14:30:45Z",
      "updated_at": "2023-01-15T14:30:45Z"
    }
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Invalid request data
  - `409 Conflict`: Username or email already exists
  - `500 Internal Server Error`: Server error

#### User Login

- **URL**: `/auth/login`
- **Method**: `POST`
- **Authentication Required**: No
- **Request Body**:
  ```json
  {
    "email": "john.doe@example.com",
    "password": "securepassword123"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john.doe@example.com",
      "created_at": "2023-01-15T14:30:45Z",
      "updated_at": "2023-01-15T14:30:45Z"
    }
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Invalid request data
  - `401 Unauthorized`: Invalid email or password
  - `500 Internal Server Error`: Server error

### Task Management

#### Create a New Task

- **URL**: `/tasks`
- **Method**: `POST`
- **Authentication Required**: Yes
- **Request Body**:
  ```json
  {
    "title": "Complete project documentation",
    "description": "Finish writing API documentation for the task manager",
    "due_date": "2023-02-15T17:00:00Z",
    "priority": "high"
  }
  ```
- **Success Response**: `201 Created`
  ```json
  {
    "id": 1,
    "user_id": 1,
    "title": "Complete project documentation",
    "description": "Finish writing API documentation for the task manager",
    "due_date": "2023-02-15T17:00:00Z",
    "priority": "high",
    "status": "todo",
    "created_at": "2023-01-20T09:15:30Z",
    "updated_at": "2023-01-20T09:15:30Z"
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Invalid request data
  - `401 Unauthorized`: Missing or invalid token
  - `500 Internal Server Error`: Server error

#### Get a Specific Task

- **URL**: `/tasks/:id`
- **Method**: `GET`
- **Authentication Required**: Yes
- **URL Parameters**: `id=[integer]` Task ID
- **Success Response**: `200 OK`
  ```json
  {
    "id": 1,
    "user_id": 1,
    "title": "Complete project documentation",
    "description": "Finish writing API documentation for the task manager",
    "due_date": "2023-02-15T17:00:00Z",
    "priority": "high",
    "status": "todo",
    "created_at": "2023-01-20T09:15:30Z",
    "updated_at": "2023-01-20T09:15:30Z"
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Invalid task ID
  - `401 Unauthorized`: Missing or invalid token
  - `404 Not Found`: Task not found
  - `500 Internal Server Error`: Server error

#### Update a Task

- **URL**: `/tasks/:id`
- **Method**: `PUT`
- **Authentication Required**: Yes
- **URL Parameters**: `id=[integer]` Task ID
- **Request Body**:
  ```json
  {
    "title": "Updated project documentation",
    "description": "Updated API documentation for the task manager",
    "due_date": "2023-02-20T17:00:00Z",
    "priority": "medium"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "id": 1,
    "user_id": 1,
    "title": "Updated project documentation",
    "description": "Updated API documentation for the task manager",
    "due_date": "2023-02-20T17:00:00Z",
    "priority": "medium",
    "status": "todo",
    "created_at": "2023-01-20T09:15:30Z",
    "updated_at": "2023-01-20T10:25:40Z"
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Invalid request data or task ID
  - `401 Unauthorized`: Missing or invalid token
  - `404 Not Found`: Task not found
  - `500 Internal Server Error`: Server error

#### Update Task Status

- **URL**: `/tasks/:id/status`
- **Method**: `PATCH`
- **Authentication Required**: Yes
- **URL Parameters**: `id=[integer]` Task ID
- **Request Body**:
  ```json
  {
    "status": "in_progress"
  }
  ```
- **Success Response**: `200 OK`
  ```json
  {
    "id": 1,
    "user_id": 1,
    "title": "Updated project documentation",
    "description": "Updated API documentation for the task manager",
    "due_date": "2023-02-20T17:00:00Z",
    "priority": "medium",
    "status": "in_progress",
    "created_at": "2023-01-20T09:15:30Z",
    "updated_at": "2023-01-21T11:30:15Z"
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Invalid request data or task ID
  - `401 Unauthorized`: Missing or invalid token
  - `404 Not Found`: Task not found
  - `500 Internal Server Error`: Server error

#### Delete a Task

- **URL**: `/tasks/:id`
- **Method**: `DELETE`
- **Authentication Required**: Yes
- **URL Parameters**: `id=[integer]` Task ID
- **Success Response**: `200 OK`
  ```json
  {
    "message": "Task deleted successfully"
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Invalid task ID
  - `401 Unauthorized`: Missing or invalid token
  - `404 Not Found`: Task not found
  - `500 Internal Server Error`: Server error

#### Get Tasks List

- **URL**: `/tasks`
- **Method**: `GET`
- **Authentication Required**: Yes
- **Query Parameters**:
  - `page=[integer]`: Page number (default: 1)
  - `page_size=[integer]`: Number of tasks per page (default: 10, max: 100)
  - `status=[string]`: Filter by status (todo, in_progress, completed)
  - `priority=[string]`: Filter by priority (low, medium, high)
  - `sort_by=[string]`: Field to sort by (created_at, due_date, priority, title)
  - `order=[string]`: Sort order (asc, desc)
- **Success Response**: `200 OK`
  ```json
  {
    "tasks": [
      {
        "id": 2,
        "user_id": 1,
        "title": "Prepare presentation",
        "description": "Create slides for project demo",
        "due_date": "2023-02-10T14:00:00Z",
        "priority": "high",
        "status": "todo",
        "created_at": "2023-01-18T13:45:20Z",
        "updated_at": "2023-01-18T13:45:20Z"
      },
      {
        "id": 1,
        "user_id": 1,
        "title": "Updated project documentation",
        "description": "Updated API documentation for the task manager",
        "due_date": "2023-02-20T17:00:00Z",
        "priority": "medium",
        "status": "in_progress",
        "created_at": "2023-01-20T09:15:30Z",
        "updated_at": "2023-01-21T11:30:15Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total_items": 2,
      "total_pages": 1
    }
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Invalid query parameters
  - `401 Unauthorized`: Missing or invalid token
  - `500 Internal Server Error`: Server error

## Health Check

- **URL**: `/health`
- **Method**: `GET`
- **Authentication Required**: No
- **Success Response**: `200 OK`
  ```json
  {
    "status": "ok"
  }
  ```

## Error Codes and Meanings

| Status Code | Description |
|-------------|-------------|
| 200 | OK - The request has succeeded |
| 201 | Created - The resource has been created |
| 400 | Bad Request - The request was invalid |
| 401 | Unauthorized - Authentication is required or failed |
| 404 | Not Found - The requested resource was not found |
| 409 | Conflict - Resource already exists (e.g., username) |
| 500 | Internal Server Error - Server encountered an error |

## Task Priority Levels

- `low`: Low priority tasks
- `medium`: Medium priority tasks (default)
- `high`: High priority tasks

## Task Status Values

- `todo`: Tasks that are not yet started (default)
- `in_progress`: Tasks that are currently being worked on
- `completed`: Tasks that are finished