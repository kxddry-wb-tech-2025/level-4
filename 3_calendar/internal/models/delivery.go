package models

import "time"

// CreateNotificationRequest is the request to create a notification
type CreateNotificationRequest struct {
	EventID   string
	Message   string
	When      time.Time
	Channel   string
	Recipient string
}

// Notification is the notification model
type Notification struct {
	EventID   string
	ID        string
	Message   string
	When      time.Time
	Channel   string
	Recipient string
}

// MessageTemplate is the template for the notification message
const (
	MessageTemplate = "You have an event %s at %s"
	NotifyBefore    = 15 * time.Minute
)

// NotifyTime calculates the time to notify for an event
func NotifyTime(start time.Time) time.Time {
	now := time.Now()
	when := start.Add(-NotifyBefore)
	if when.Before(now) {
		return now
	}

	return when
}
