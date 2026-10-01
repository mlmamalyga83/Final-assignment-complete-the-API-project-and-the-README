package models

import "time"

// User model
// @Description Represents a user account
type User struct {
	ID        int       `json:"id" example:"1"`
	Email     string    `json:"email" example:"user@example.com"`
	Name      string    `json:"name" example:"John Doe"`
	CreatedAt time.Time `json:"created_at" swaggertype:"string" format:"date-time"`
}

// CreateUserRequest model
// @Description Request body for creating a new user
type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
	Name  string `json:"name" validate:"required,max=100" example:"John Doe"`
}

// UpdateUserRequest model
// @Description Request body for updating an existing user (all fields optional)
type UpdateUserRequest struct {
	Email *string `json:"email" validate:"omitempty,email" example:"newemail@example.com"`
	Name  *string `json:"name" validate:"omitempty,max=100" example:"Jane Doe"`
}

// ErrorResponse model
// @Description Standard error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse model
// @Description Standard success response wrapper
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}
