package inmemory

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sort"
	"sync"
	"time"
)

type EventRepository struct {
	mu     sync.RWMutex
	events map[int64][]*domain.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		events: make(map[int64][]*domain.Event),
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

	_, ok := r.events[event.SensorID]
	if !ok {
		r.events[event.SensorID] = []*domain.Event{}
	}
	r.events[event.SensorID] = r.insertEvent(r.events[event.SensorID], event)

	return nil
}

func (r *EventRepository) insertEvent(events []*domain.Event, event *domain.Event) []*domain.Event {
	//if len(events) == 0 {
	//	events = []*domain.Event{}
	//	events = append(events, event)
	//	return events
	//}
	ind := sort.Search(len(events), func(i int) bool {
		return !(events)[i].Timestamp.Before(event.Timestamp)
	})
	events = append(events, &domain.Event{})
	copy((events)[ind+1:], (events)[ind:])
	events[ind] = event
	return events
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	events, ok := r.events[id]
	if !ok {
		return nil, fmt.Errorf("no events yet on sensor %d: %w", id, usecase.ErrEventNotFound)
	}
	return events[len(events)-1], nil
}

func (r *EventRepository) GetEventsInRangeBySensorID(ctx context.Context, id int64, from, to time.Time) ([]*domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	events, ok := r.events[id]
	if !ok {
		return nil, fmt.Errorf("no events yet on sensor %d: %w", id, usecase.ErrEventNotFound)
	}

	start := sort.Search(len(events)-1, func(i int) bool {
		return !(events)[i].Timestamp.Before(from)
	})
	end := sort.Search(len(events)-1, func(i int) bool {
		return (events)[i].Timestamp.After(to)
	}) + 1

	return (events)[start:end], nil
}
