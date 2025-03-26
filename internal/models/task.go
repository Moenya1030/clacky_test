package models

import (
	"time"

	"gorm.io/gorm"
)

// Priority represents the priority level of a task
type Priority string

// Status represents the current state of a task
type Status string

const (
	// Priority levels
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"

	// Status options
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
)

// Task represents the task model in the database
type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	DueDate     *time.Time     `json:"due_date"`
	Priority    Priority       `gorm:"type:enum('low','medium','high');default:'medium'" json:"priority"`
	Status      Status         `gorm:"type:enum('todo','in_progress','completed');default:'todo'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for the Task model
func (Task) TableName() string {
	return "tasks"
}