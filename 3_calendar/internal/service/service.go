package service

import (
	"calendar/internal/models"
	"context"
)

// NotificationCreator is the interface for creating a notification.
type NotificationCreator interface {
	CreateNotification(ctx context.Context, notification models.CreateNotificationRequest) (string, error)
}

// NotificationDeleter is the interface for deleting a notification.
type NotificationDeleter interface {
	DeleteNotificationByID(ctx context.Context, id string) error
}

// NotificationGetter is the interface for getting a notification.
type NotificationGetter interface {
	GetNotificationByID(ctx context.Context, id string) (models.Notification, error)
	GetNotificationIDsByEventID(ctx context.Context, eventID string) ([]string, error)
}

// NotificationAdder is the interface for adding a notification.
type NotificationAdder interface {
	AddNotification(ctx context.Context, notification models.CreateNotificationRequest) error
}

// EventCreator is the interface for creating an event.
type EventCreator interface {
	CreateEvent(ctx context.Context, event models.CreateEventRequest) (string, error)
}

// EventCreatorID is the interface for creating an event with an ID.
type EventCreatorID interface {
	CreateEventWithID(ctx context.Context, id string, event models.CreateEventRequest) error
}

// EventGetter is the interface for getting an event.
type EventGetter interface {
	GetEvent(ctx context.Context, id string) (models.Event, error)
	GetAllEvents(ctx context.Context) ([]models.Event, error)
}

// EventUpdater is the interface for updating an event.
type EventUpdater interface {
	UpdateEvent(ctx context.Context, id string, event models.UpdateEventRequest) error
}

// EventDeleter is the interface for deleting an event.
type EventDeleter interface {
	DeleteEvent(ctx context.Context, id string) error
}
