package inmemory

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
)

type EventRepository struct {
	mu     sync.RWMutex
	events map[int64]*domain.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		events: make(map[int64]*domain.Event),
	}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if event == nil {
		return fmt.Errorf("event is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	lastEvent, ok := r.events[event.SensorID]
	if ok && lastEvent.Timestamp.After(event.Timestamp) {
		return nil
	}

	r.events[event.SensorID] = event

	return nil
}

// GetLastEventBySensorID ; то же самое про ошибку из usecase, что и в sensor.GetSensorByID
func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	event, ok := r.events[id]
	var myErr error
	if !ok {
		return nil, fmt.Errorf("no events yet on sensor %d: %w", id, usecase.ErrEventNotFound)
	}
	return event, myErr
}
