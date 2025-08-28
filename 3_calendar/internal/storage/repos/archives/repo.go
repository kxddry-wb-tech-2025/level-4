package archives

import (
	"calendar/internal/models"
	"calendar/internal/models/log"
	"calendar/internal/storage"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool    *pgxpool.Pool
	logs    chan<- log.Entry
	mainCtx context.Context
}

func NewRepository(ctx context.Context, pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, mainCtx: ctx}
}

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

func (r *Repository) Archive(ctx context.Context, event models.Event) error {
	query := `
	INSERT INTO archives (event_id, title, description, start, end, notify, email)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var err error
	if tx, ok := storage.GetTx(ctx); ok {
		_, err = tx.Exec(ctx, query, event.ID, event.Title, event.Description, event.Start, event.End, event.Notify, event.Email)
	} else {
		r.sendLog(log.Warn("no tx found, using pool", map[string]any{
			"op":         "archive",
			"event_id":   event.ID,
			"repository": "archives",
		}))
		_, err = r.pool.Exec(ctx, query, event.ID, event.Title, event.Description, event.Start, event.End, event.Notify, event.Email)
	}

	return err
}
