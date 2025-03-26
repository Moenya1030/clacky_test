package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"task-manager/internal/models"
	"task-manager/pkg/database"
)

// TaskRequest defines the data needed to create or update a task
type TaskRequest struct {
	Title       string
	Description string
	DueDate     *time.Time
	Priority    models.Priority
	UserID      uint
}

// TaskStatusRequest defines the data needed to update a task status
type TaskStatusRequest struct {
	Status models.Status
	UserID uint
}

// TaskFilterOptions defines the options for filtering and sorting tasks
type TaskFilterOptions struct {
	UserID   uint
	Status   string
	Priority string
	SortBy   string
	Order    string
	Page     int
	PageSize int
}

// PaginatedTasksResponse represents a paginated list of tasks
type PaginatedTasksResponse struct {
	Tasks       []models.Task
	CurrentPage int
	PageSize    int
	TotalItems  int64
	TotalPages  int64
}

// TaskService provides methods for task-related operations
type TaskService struct {
	db *gorm.DB
}

// NewTaskService creates a new instance of TaskService
func NewTaskService() *TaskService {
	return &TaskService{
		db: database.GetDB(),
	}
}

// CreateTask creates a new task for the user
func (s *TaskService) CreateTask(req TaskRequest) (*models.Task, error) {
	task := models.Task{
		UserID:      req.UserID,
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
	if err := s.db.Create(&task).Error; err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return &task, nil
}

// GetTaskByID retrieves a task by ID if it belongs to the specified user
func (s *TaskService) GetTaskByID(taskID uint, userID uint) (*models.Task, error) {
	var task models.Task
	result := s.db.Where("id = ? AND user_id = ?", taskID, userID).First(&task)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, fmt.Errorf("failed to retrieve task: %w", result.Error)
	}
	return &task, nil
}

// UpdateTask updates an existing task if it belongs to the specified user
func (s *TaskService) UpdateTask(taskID uint, req TaskRequest) (*models.Task, error) {
	// Find task by ID and ensure it belongs to the user
	task, err := s.GetTaskByID(taskID, req.UserID)
	if err != nil {
		return nil, err
	}

	// Update task fields
	task.Title = req.Title
	task.Description = req.Description
	task.DueDate = req.DueDate
	if req.Priority != "" {
		task.Priority = req.Priority
	}

	// Save updated task
	if err := s.db.Save(task).Error; err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// UpdateTaskStatus updates only the status of a task
func (s *TaskService) UpdateTaskStatus(taskID uint, req TaskStatusRequest) (*models.Task, error) {
	// Find task by ID and ensure it belongs to the user
	task, err := s.GetTaskByID(taskID, req.UserID)
	if err != nil {
		return nil, err
	}

	// Update task status
	task.Status = req.Status

	// Save updated task
	if err := s.db.Save(task).Error; err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	return task, nil
}

// DeleteTask deletes a task if it belongs to the specified user
func (s *TaskService) DeleteTask(taskID uint, userID uint) error {
	// Find task by ID and ensure it belongs to the user
	task, err := s.GetTaskByID(taskID, userID)
	if err != nil {
		return err
	}

	// Delete the task (soft delete with GORM)
	if err := s.db.Delete(task).Error; err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// GetTasks retrieves tasks with pagination, filtering, and sorting
func (s *TaskService) GetTasks(options TaskFilterOptions) (*PaginatedTasksResponse, error) {
	// Set default pagination values if not provided
	page := options.Page
	if page < 1 {
		page = 1
	}

	pageSize := options.PageSize
	if pageSize < 1 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Start building the query
	query := s.db.Model(&models.Task{}).Where("user_id = ?", options.UserID)

	// Apply filters if provided
	if options.Status != "" {
		query = query.Where("status = ?", options.Status)
	}
	if options.Priority != "" {
		query = query.Where("priority = ?", options.Priority)
	}

	// Determine sorting
	sortBy := "created_at" // default sort field
	if options.SortBy != "" {
		sortBy = options.SortBy
	}

	order := "desc" // default order
	if options.Order != "" {
		order = options.Order
	}

	// Get total count of matching tasks
	var totalTasks int64
	if err := query.Count(&totalTasks).Error; err != nil {
		return nil, fmt.Errorf("failed to count tasks: %w", err)
	}

	// Apply sorting, pagination, and execute query
	var tasks []models.Task
	if err := query.Order(sortBy + " " + order).
		Limit(pageSize).
		Offset(offset).
		Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}

	// Calculate total pages
	totalPages := (totalTasks + int64(pageSize) - 1) / int64(pageSize)

	// Return response with pagination metadata
	return &PaginatedTasksResponse{
		Tasks:       tasks,
		CurrentPage: page,
		PageSize:    pageSize,
		TotalItems:  totalTasks,
		TotalPages:  totalPages,
	}, nil
}