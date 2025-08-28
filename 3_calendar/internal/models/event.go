package models

import "time"

// Event is the event model
type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Notify      bool      `json:"notify"`
	Email       string    `json:"email"`
}

// CreateEventRequest is the request to create an event
type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Notify      bool      `json:"notify"`
	Email       string    `json:"email"`
}

// UpdateEventRequest is the request to update an event
type UpdateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Notify      bool      `json:"notify"`
	Email       string    `json:"email"`
}
