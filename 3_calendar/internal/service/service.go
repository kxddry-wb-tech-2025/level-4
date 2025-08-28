package service

import (
	"calendar/internal/models"
	"context"
)

type NotificationCreator interface {
	CreateNotification(ctx context.Context, notification models.CreateNotificationRequest) (string, error)
}

type NotificationDeleter interface {
	DeleteNotificationByID(ctx context.Context, id string) error
	DeleteAllNotificationsByEventID(ctx context.Context, eventID string) error
}

type NotificationGetter interface {
	GetNotificationByID(ctx context.Context, id string) (models.Notification, error)
	GetNotificationIDsByEventID(ctx context.Context, eventID string) ([]string, error)
}

type NotificationAdder interface {
	AddNotification(ctx context.Context, notification models.CreateNotificationRequest) error
}

type EventCreator interface {
	CreateEvent(ctx context.Context, event models.CreateEventRequest) (string, error)
}

type EventGetter interface {
	GetEvent(ctx context.Context, id string) (models.Event, error)
	GetAllEvents(ctx context.Context) ([]models.Event, error)
}

type EventUpdater interface {
	UpdateEvent(ctx context.Context, id string, event models.UpdateEventRequest) error
}

type EventDeleter interface {
	DeleteEvent(ctx context.Context, id string) error
}
