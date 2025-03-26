package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"task-manager/internal/middlewares"
	"task-manager/internal/models"
	"task-manager/pkg/database"
)

// TaskRequest represents the request body for creating/updating a task
type TaskRequest struct {
	Title       string        `json:"title" binding:"required,max=200"`
	Description string        `json:"description"`
	DueDate     *time.Time    `json:"due_date"`
	Priority    models.Priority `json:"priority" binding:"omitempty,oneof=low medium high"`
}

// TaskStatusRequest represents the request body for updating task status
type TaskStatusRequest struct {
	Status models.Status `json:"status" binding:"required,oneof=todo in_progress completed"`
}

// PaginationQuery represents the query parameters for pagination
type PaginationQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=5,max=100"`
}

// TaskFilterQuery represents the query parameters for filtering tasks
type TaskFilterQuery struct {
	Status   string `form:"status" binding:"omitempty,oneof=todo in_progress completed"`
	Priority string `form:"priority" binding:"omitempty,oneof=low medium high"`
	SortBy   string `form:"sort_by" binding:"omitempty,oneof=created_at due_date priority title"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

// CreateTask handles the creation of a new task
func CreateTask(c *gin.Context) {
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Create new task
	task := models.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      models.StatusTodo, // Default status is todo
	}

	// Set priority if provided, otherwise use default (medium)
	if req.Priority != "" {
		task.Priority = req.Priority
	} else {
		task.Priority = models.PriorityMedium
	}

	// Save task to database
	if err := database.GetDB().Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetTask retrieves a single task by its ID
func GetTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Find task by ID and ensure it belongs to the authenticated user
	var task models.Task
	result := database.GetDB().Where("id = ? AND user_id = ?", taskID, userID).First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve task: " + result.Error.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask updates a task's details
func UpdateTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	// Parse request body
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Find task by ID and ensure it belongs to the authenticated user
	var task models.Task
	result := database.GetDB().Where("id = ? AND user_id = ?", taskID, userID).First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve task: " + result.Error.Error(),
			})
		}
		return
	}

	// Update task fields
	task.Title = req.Title
	task.Description = req.Description
	task.DueDate = req.DueDate
	if req.Priority != "" {
		task.Priority = req.Priority
	}

	// Save updated task
	if err := database.GetDB().Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTaskStatus updates only the status of a task
func UpdateTaskStatus(c *gin.Context) {
	// Get task ID from URL parameter
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	// Parse request body
	var req TaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Find task by ID and ensure it belongs to the authenticated user
	var task models.Task
	result := database.GetDB().Where("id = ? AND user_id = ?", taskID, userID).First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve task: " + result.Error.Error(),
			})
		}
		return
	}

	// Update task status
	task.Status = req.Status

	// Save updated task
	if err := database.GetDB().Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask deletes a task by its ID
func DeleteTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	// Get user ID from context
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Check if task exists and belongs to the user
	var task models.Task
	result := database.GetDB().Where("id = ? AND user_id = ?", taskID, userID).First(&task)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve task: " + result.Error.Error(),
			})
		}
		return
	}

	// Delete the task (soft delete with GORM)
	if err := database.GetDB().Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}

// GetTasks retrieves a list of tasks with pagination, filtering, and sorting
func GetTasks(c *gin.Context) {
	// Get user ID from context
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// Parse pagination parameters
	var pagination PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid pagination parameters: " + err.Error(),
		})
		return
	}

	// Set default pagination values if not provided
	page := pagination.Page
	if page == 0 {
		page = 1
	}

	pageSize := pagination.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Parse filter parameters
	var filter TaskFilterQuery
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid filter parameters: " + err.Error(),
		})
		return
	}

	// Start building the query
	query := database.GetDB().Model(&models.Task{}).Where("user_id = ?", userID)

	// Apply filters if provided
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}

	// Determine sorting
	sortBy := "created_at" // default sort field
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}

	order := "desc" // default order
	if filter.Order != "" {
		order = filter.Order
	}

	// Get total count of matching tasks
	var totalTasks int64
	if err := query.Count(&totalTasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count tasks: " + err.Error(),
		})
		return
	}

	// Apply sorting, pagination, and execute query
	var tasks []models.Task
	if err := query.Order(sortBy + " " + order).
		Limit(pageSize).
		Offset(offset).
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve tasks: " + err.Error(),
		})
		return
	}

	// Calculate total pages
	totalPages := (totalTasks + int64(pageSize) - 1) / int64(pageSize)

	// Return response with pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
		"pagination": gin.H{
			"current_page": page,
			"page_size":    pageSize,
			"total_items":  totalTasks,
			"total_pages":  totalPages,
		},
	})
}