package models

import "time"

// Event is the event model
type Event struct {
	ID          string    `json:"id" validate:"required,uuid"`
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Description string    `json:"description" validate:"required,min=1,max=1024"`
	Start       time.Time `json:"start" validate:"required,future"`
	End         time.Time `json:"end" validate:"required,future"`
	Notify      bool      `json:"notify" validate:"omitempty,boolean"`
	Email       string    `json:"email" validate:"omitempty,email"`
}

// CreateEventRequest is the request to create an event
type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Description string    `json:"description" validate:"required,min=1,max=1024"`
	Start       time.Time `json:"start" validate:"required,future"`
	End         time.Time `json:"end" validate:"required,future"`
	Notify      bool      `json:"notify" validate:"omitempty,boolean"`
	Email       string    `json:"email" validate:"omitempty,email"`
}

// UpdateEventRequest is the request to update an event
type UpdateEventRequest struct {
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Description string    `json:"description" validate:"required,min=1,max=1024"`
	Start       time.Time `json:"start" validate:"required,future"`
	End         time.Time `json:"end" validate:"required,future"`
	Notify      bool      `json:"notify" validate:"omitempty,boolean"`
	Email       string    `json:"email" validate:"omitempty,email"`
}
