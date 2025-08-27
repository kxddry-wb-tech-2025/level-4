package events

import (
	"calendar/internal/models"
	"calendar/internal/storage"
	"context"
	"errors"
)

type EventRepository interface {
	Create(event models.CreateEventRequest) (string, error)
	GetAll() ([]models.Event, error)
	Get(id string) (models.Event, error)
	Update(id string, event models.UpdateEventRequest) error
	Delete(id string) error
}

type Service struct {
	repo EventRepository
}

func NewService(repo EventRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEvent(ctx context.Context, event models.CreateEventRequest) (string, error) {
	id, err := s.repo.Create(event)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) GetEvents(ctx context.Context) ([]models.Event, error) {
	events, err := s.repo.GetAll()
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return []models.Event{}, nil
		}

		return nil, err
	}

	return events, nil
}

func (s *Service) GetEvent(ctx context.Context, id string) (models.Event, error) {
	event, err := s.repo.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.Event{}, models.ErrNotFound
		}

		return models.Event{}, err
	}

	return event, nil
}

func (s *Service) UpdateEvent(ctx context.Context, id string, event models.UpdateEventRequest) error {
	err := s.repo.Update(id, event)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.ErrNotFound
		}

		return err
	}

	return nil
}

func (s *Service) DeleteEvent(ctx context.Context, id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return models.ErrNotFound
		}

		return err
	}

	return nil
}
