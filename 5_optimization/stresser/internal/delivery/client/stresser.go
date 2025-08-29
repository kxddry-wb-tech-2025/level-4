package client

import (
	"context"
	"errors"
	"stresser/internal/models"
	"stresser/internal/models/log"

	"github.com/go-resty/resty/v2"
)

// Stresser is a client that can stress a victim server.
type Stresser struct {
	logs   chan log.Entry
	ctx    context.Context
	client *resty.Client
	host   string
}

// sendLog sends a log entry to the logs channel.
func (s *Stresser) sendLog(entry log.Entry) {
	if s.logs == nil {
		return
	}

	if entry.Level >= log.LevelWarn {
		go func() {
			select {
			case s.logs <- entry:
			case <-s.ctx.Done():
			}
		}()
	} else {
		select {
		case s.logs <- entry:
		case <-s.ctx.Done():
		default:
		}
	}
}

// Logs returns a channel to receive log entries.
func (s *Stresser) Logs() <-chan log.Entry {
	logs := make(chan log.Entry, log.ChannelCapacity)
	s.logs = logs

	// Start a goroutine to close the channel when context is canceled
	go func() {
		<-s.ctx.Done()
		close(s.logs)
	}()

	return logs
}

// NewStresser creates a new Stresser.
func NewStresser(victim string, ctx context.Context) *Stresser {
	return &Stresser{client: resty.New(), host: victim, ctx: ctx}
}

// Stress sends orders to the victim server.
func (s *Stresser) Stress(orders []models.Order, reuse bool) {
	process := func() {
		for _, order := range orders {
			select {
			case <-s.ctx.Done():
				return
			default:
				if err := s.sendOrder(order); err != nil {
					s.sendLog(log.Error(err, "failed to send order", map[string]any{"order": order}))
				}
			}
		}
	}

	if !reuse {
		go process()
	} else {
		go func() {
			for {
				select {
				case <-s.ctx.Done():
					return
				default:
					process()
				}
			}
		}()
	}
}

// sendOrder sends an order to the victim server.
func (s *Stresser) sendOrder(order models.Order) error {
	resp, err := s.client.R().SetBody(order).SetHeader("Content-Type", "application/json").Post(s.host + "/order")
	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(resp.Status())
	}

	return nil
}
