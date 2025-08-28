package notifications

import (
	"calendar/internal/models"
	"calendar/internal/models/log"
	"calendar/internal/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// Repository provides DB access for notifications
type Repository struct {
	pool    *pgxpool.Pool
	logs    chan<- log.Entry
	mainCtx context.Context
}

// NewRepository creates a new repository
func NewRepository(ctx context.Context, pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, mainCtx: ctx}
}

// Logs returns the channel for the logs
func (r *Repository) Logs() <-chan log.Entry {
	logs := make(chan log.Entry, log.ChannelCapacity)
	r.logs = logs
	return logs
}

func (r *Repository) sendLog(entry log.Entry) {
	if r.logs == nil {
		return
	}
	if entry.Level >= log.LevelWarn {
		go func() {
			select {
			case r.logs <- entry:
			case <-r.mainCtx.Done():
			}
		}()
	} else {
		select {
		case r.logs <- entry:
		default:
		}
	}
}

// Create creates a new notification
func (r *Repository) Create(ctx context.Context, notification models.CreateNotificationRequest) (string, error) {
	query := `
	INSERT INTO notifications (event_id, message, "when", channel, recipient)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	var row pgx.Row
	if tx, ok := storage.GetTx(ctx); ok {
		row = tx.QueryRow(ctx, query, notification.EventID, notification.Message, notification.When, notification.Channel, notification.Recipient)
	} else {
		r.sendLog(log.Warn("no tx found, using pool", echo.Map{
			"op":   "Create",
			"repo": "notifications",
		}))
		row = r.pool.QueryRow(ctx, query, notification.EventID, notification.Message, notification.When, notification.Channel, notification.Recipient)
	}
	var id string
	err := row.Scan(&id)
	return id, err
}

// GetIDsByEventID gets all notification ids by event id
func (r *Repository) GetIDsByEventID(ctx context.Context, eventID string) ([]string, error) {
	query := `
	SELECT id
	FROM notifications
	WHERE event_id = $1
	`
	var rows pgx.Rows
	var err error
	if tx, ok := storage.GetTx(ctx); ok {
		rows, err = tx.Query(ctx, query, eventID)
	} else {
		r.sendLog(log.Warn("no tx found, using pool", echo.Map{
			"op":   "GetIDsByEventID",
			"repo": "notifications",
		}))
		rows, err = r.pool.Query(ctx, query, eventID)
	}
	if err != nil {
		return nil, err
	}
	var ids []string
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []string{}, storage.ErrNotFound
		}
		return nil, err
	}
	return ids, nil
}

// DeleteAllByEventID deletes all notifications by event id
func (r *Repository) DeleteAllByEventID(ctx context.Context, eventID string) error {
	query := `
	DELETE FROM notifications
	WHERE event_id = $1
	`
	var err error
	if tx, ok := storage.GetTx(ctx); ok {
		_, err = tx.Exec(ctx, query, eventID)
	} else {
		r.sendLog(log.Warn("no tx found, using pool", echo.Map{
			"op":   "DeleteAllByEventID",
			"repo": "notifications",
		}))
		_, err = r.pool.Exec(ctx, query, eventID)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}
	return err
}

// GetByID gets a notification by id
func (r *Repository) GetByID(ctx context.Context, id string) (models.Notification, error) {
	query := `
	SELECT id, event_id, message, "when", channel, recipient
	FROM notifications
	WHERE id = $1
	`
	var row pgx.Row
	if tx, ok := storage.GetTx(ctx); ok {
		row = tx.QueryRow(ctx, query, id)
	} else {
		r.sendLog(log.Warn("no tx found, using pool", echo.Map{
			"op":   "GetByID",
			"repo": "notifications",
		}))
		row = r.pool.QueryRow(ctx, query, id)
	}
	var n models.Notification
	if err := row.Scan(&n.ID, &n.EventID, &n.Message, &n.When, &n.Channel, &n.Recipient); err != nil {
		return models.Notification{}, err
	}
	return n, nil
}

// DeleteByID deletes a notification by id
func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	query := `
	DELETE FROM notifications
	WHERE id = $1
	`
	var err error
	if tx, ok := storage.GetTx(ctx); ok {
		_, err = tx.Exec(ctx, query, id)
	} else {
		r.sendLog(log.Warn("no tx found, using pool", echo.Map{
			"op":   "DeleteByID",
			"repo": "notifications",
		}))
		_, err = r.pool.Exec(ctx, query, id)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}
	return nil
}
