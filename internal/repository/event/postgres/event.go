package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/domain"
	"time"
)

var ErrEventNotFound = errors.New("event not found")

type EventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		pool,
	}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	// TODO добавьте реализацию
	return nil
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	// TODO добавьте реализацию
	return nil, nil
}

func (r *EventRepository) GetEventsInRangeBySensorID(ctx context.Context, id int64, from, to time.Time) ([]*domain.Event, error) {
	//TODO implement me
	panic("implement me")
}
