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
	if e.eventRepo == nil && e.sensorRepo == nil {
		return ErrInvalidEventTimestamp
	}
	if e.sensorRepo == nil {
		return ErrSensorNotFound
	}

	sensor, err := e.sensorRepo.GetSensorBySerialNumber(ctx, event.SensorSerialNumber)
	if err != nil {
		return err
	}

	event.SensorID = sensor.ID
	sensor.CurrentState = event.Payload
	sensor.LastActivity = event.Timestamp

	if err := e.eventRepo.SaveEvent(ctx, event); err != nil {
		return err
	}

	if err := e.sensorRepo.SaveSensor(ctx, sensor); err != nil {
		return err
	}

	return nil
}

func (e *Event) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	event, err := e.eventRepo.GetLastEventBySensorID(ctx, id)

	return event, err
}
