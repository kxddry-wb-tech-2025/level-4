package models

import "time"

type CreateNotificationRequest struct {
	EventID   string
	Message   string
	When      time.Time
	Channel   string
	Recipient string
}

type DeleteNotificationsRequest struct {
	EventID string
}

type Notification struct {
	EventID   string
	ID        string
	Message   string
	When      time.Time
	Channel   string
	Recipient string
}
