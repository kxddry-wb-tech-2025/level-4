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

const (
	MessageTemplate = "You have an event %s at %s"
	NotifyBefore    = 15 * time.Minute
)

func NotifyTime(start time.Time) time.Time {
	now := time.Now()
	when := start.Add(-NotifyBefore)
	if when.Before(now) {
		return now
	}

	return when
}
