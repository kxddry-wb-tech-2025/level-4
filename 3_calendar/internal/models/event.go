package models

import "time"

type Event struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Notify      bool      `json:"notify"`
	Email       string    `json:"email"`
}

type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Notify      bool      `json:"notify"`
	Email       string    `json:"email"`
}

type UpdateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Notify      bool      `json:"notify"`
	Email       string    `json:"email"`
}
