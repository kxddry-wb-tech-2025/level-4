package models

import "errors"

// Errors
var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidEvent = errors.New("invalid event")
)
