package events

import (
	"calendar/internal/models"
	"calendar/internal/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository is a repository for events
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new repository
func NewRepository(ctx context.Context, connStr string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return &Repository{pool: pool}, pool.Ping(ctx)
}

// Create creates a new event
func (r *Repository) Create(ctx context.Context, event models.CreateEventRequest) (string, error) {
	query := `
	INSERT INTO events (title, description, start, end, notify, email)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`
	row := r.pool.QueryRow(ctx, query, event.Title, event.Description, event.Start, event.End, event.Notify, event.Email)
	var id string
	err := row.Scan(&id)
	return id, err
}

// GetAll gets all events
func (r *Repository) GetAll(ctx context.Context) ([]models.Event, error) {
	query := `
	SELECT id, title, description, start, end, notify, email
	FROM events
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Start, &event.End, &event.Notify, &event.Email)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.Event{}, storage.ErrNotFound
		}
		return nil, err
	}
	return events, nil
}

// Get gets an event by id
func (r *Repository) Get(ctx context.Context, id string) (models.Event, error) {
	query := `
	SELECT id, title, description, start, end, notify, email
	FROM events
	WHERE id = $1
	`
	row := r.pool.QueryRow(ctx, query, id)
	var event models.Event
	err := row.Scan(&event.ID, &event.Title, &event.Description, &event.Start, &event.End, &event.Notify, &event.Email)
	if err != nil {
		return models.Event{}, err
	}
	return event, nil
}

// Update updates an event
func (r *Repository) Update(ctx context.Context, id string, event models.UpdateEventRequest) error {
	query := `
	UPDATE events
	SET title = $1, description = $2, start = $3, end = $4, notify = $5, email = $6
	WHERE id = $7
	`
	_, err := r.pool.Exec(ctx, query, event.Title, event.Description, event.Start, event.End, event.Notify, event.Email, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}
	return nil
}

// Delete deletes an event
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM events
	WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrNotFound
		}
		return err
	}
	return nil
}
