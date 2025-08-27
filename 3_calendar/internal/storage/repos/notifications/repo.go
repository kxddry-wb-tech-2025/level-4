package notifications

import (
	"calendar/internal/models"
	"calendar/internal/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository interface {
	Create(notification models.Notification) error
	GetIDsByEventID(eventID string) ([]string, error)
	DeleteAllByEventID(eventID string) error
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &Repository{pool: pool}, pool.Ping(ctx)
}

func (r *Repository) Create(ctx context.Context, notification models.CreateNotificationRequest) (string, error) {
	query := `
	INSERT INTO notifications (event_id, message, when, channel, recipient)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	row := r.pool.QueryRow(ctx, query, notification.EventID, notification.Message, notification.When, notification.Channel, notification.Recipient)
	var id string
	err := row.Scan(&id)
	return id, err
}

func (r *Repository) GetIDsByEventID(ctx context.Context, eventID string) ([]string, error) {
	query := `
	SELECT id
	FROM notifications
	WHERE event_id = $1
	`
	rows, err := r.pool.Query(ctx, query, eventID)
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

func (r *Repository) DeleteAllByEventID(ctx context.Context, eventID string) error {
	query := `
	DELETE FROM notifications
	WHERE event_id = $1
	`
	_, err := r.pool.Exec(ctx, query, eventID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}
	return err
}
