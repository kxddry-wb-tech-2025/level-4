package txmanager

import (
	"calendar/internal/models"
	"context"
)

func (t *tx) CreateEvent(ctx context.Context, event models.CreateEventRequest) (string, error) {
	return t.repos.Events.Create(ctx, event)
}

func (t *tx) CreateEventWithID(ctx context.Context, id string, event models.CreateEventRequest) error {
	return t.repos.Events.CreateWithID(ctx, id, event)
}

func (t *tx) GetAllEvents(ctx context.Context) ([]models.Event, error) {
	return t.repos.Events.GetAll(ctx)
}

func (t *tx) GetEvent(ctx context.Context, id string) (models.Event, error) {
	return t.repos.Events.Get(ctx, id)
}

func (t *tx) UpdateEvent(ctx context.Context, id string, event models.UpdateEventRequest) error {
	return t.repos.Events.Update(ctx, id, event)
}

func (t *tx) DeleteEvent(ctx context.Context, id string) error {
	return t.repos.Events.Delete(ctx, id)
}

func (t *tx) ArchiveEvent(ctx context.Context, event models.Event) error {
	return t.repos.Archives.Archive(ctx, event)
}

func (t *tx) CreateNotification(ctx context.Context, notification models.CreateNotificationRequest) (string, error) {
	return t.repos.Notifications.Create(ctx, notification)
}

func (t *tx) GetNotificationIDsByEventID(ctx context.Context, eventID string) ([]string, error) {
	return t.repos.Notifications.GetIDsByEventID(ctx, eventID)
}

func (t *tx) DeleteAllNotificationsByEventID(ctx context.Context, eventID string) error {
	return t.repos.Notifications.DeleteAllByEventID(ctx, eventID)
}

func (t *tx) DeleteNotificationByID(ctx context.Context, id string) error {
	return t.repos.Notifications.DeleteByID(ctx, id)
}

func (t *tx) GetNotificationByID(ctx context.Context, id string) (models.Notification, error) {
	return t.repos.Notifications.GetByID(ctx, id)
}
