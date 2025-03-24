package usecase

import (
	"context"
	"homework/internal/domain"
)

type Event struct {
	eventRepo  EventRepository
	sensorRepo SensorRepository
}

func NewEvent(er EventRepository, sr SensorRepository) *Event {
	return &Event{
		eventRepo:  er,
		sensorRepo: sr,
	}
}

func (e *Event) ReceiveEvent(ctx context.Context, event *domain.Event) error {
	if event == nil {
		return ErrEventNotFound
	}
	if e.sensorRepo == nil {
		return ErrSensorNotFound
	}

	if _, err := e.sensorRepo.GetSensorByID(ctx, event.SensorID); err != nil {
		return err
	}

	if e.eventRepo == nil {
		return ErrInvalidEventTimestamp
	}

	return e.eventRepo.SaveEvent(ctx, event)
}

func (e *Event) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	event, err := e.eventRepo.GetLastEventBySensorID(ctx, id)

	return event, err
}
