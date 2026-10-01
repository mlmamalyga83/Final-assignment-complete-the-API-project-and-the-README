package models

import "time"

// Todo model
// @Description Represents a single task with its metadata
type Todo struct {
	ID          int       `json:"id" example:"1"`
	UserID      int       `json:"user_id" example:"1"`
	Title       string    `json:"title" example:"Buy groceries"`
	Description string    `json:"description" example:"Milk, bread, eggs, cheese"`
	Done        bool      `json:"done" example:"false"`
	CreatedAt   time.Time `json:"created_at" swaggertype:"string" format:"date-time"`
	UpdatedAt   time.Time `json:"updated_at" swaggertype:"string" format:"date-time"`
}

// CreateTodoRequest model
// @Description Request body for creating a new task
type CreateTodoRequest struct {
	UserID      int    `json:"user_id" validate:"required" example:"1"`
	Title       string `json:"title" validate:"required,max=200" example:"Buy groceries"`
	Description string `json:"description" validate:"max=1000" example:"Milk, bread, eggs, cheese"`
	Done        bool   `json:"done" example:"false"`
}

// UpdateTodoRequest model
// @Description Request body for updating an existing task (all fields optional)
type UpdateTodoRequest struct {
	Title       *string `json:"title" validate:"omitempty,max=200" example:"Buy groceries and snacks"`
	Description *string `json:"description" validate:"omitempty,max=1000" example:"Updated description"`
	Done        *bool   `json:"done" example:"true"`
}
