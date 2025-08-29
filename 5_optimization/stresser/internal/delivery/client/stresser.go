package client

import (
	"context"
	"errors"
	"math/rand"
	"stresser/internal/models"
	"stresser/internal/models/log"
	"sync"
	"time"

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
func (s *Stresser) Stress(orders []models.Order, reuse bool, concurrency int, operation string) {
	process := func() {
		for _, order := range orders {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)

			select {
			case <-s.ctx.Done():
				return
			default:
				switch operation {
				case "create":
					if err := s.sendOrder(order); err != nil {
						s.sendLog(log.Error(err, "failed to send order", map[string]any{"order": order}))
					}
				case "get":
					if err := s.getOrder(order.OrderUID); err != nil {
						s.sendLog(log.Error(err, "failed to get order", map[string]any{"order": order}))
					}
				}
			}
		}
	}

	wg := new(sync.WaitGroup)

	for range concurrency {
		wg.Go(func() {
			if !reuse {
				process()
			} else {
				for {
					select {
					case <-s.ctx.Done():
						return
					default:
						process()
					}
				}
			}
		})
	}

	wg.Wait()
}

func (s *Stresser) getOrder(orderUID string) error {
	resp, err := s.client.R().Get(s.host + "/order/" + orderUID)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New(resp.String())
	}

	return nil
}

// sendOrder sends an order to the victim server.
func (s *Stresser) sendOrder(order models.Order) error {
	resp, err := s.client.R().SetBody(order).SetHeader("Content-Type", "application/json").Post(s.host + "/order")
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New(resp.String())
	}

	return nil
}
