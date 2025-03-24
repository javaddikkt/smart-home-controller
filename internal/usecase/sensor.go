package usecase

import (
	"context"
	"homework/internal/domain"
)

type Sensor struct {
	sensorRepo SensorRepository
}

func NewSensor(sr SensorRepository) *Sensor {
	return &Sensor{
		sensorRepo: sr,
	}
}

func (s *Sensor) RegisterSensor(ctx context.Context, sensor *domain.Sensor) (*domain.Sensor, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return sensor, s.sensorRepo.SaveSensor(ctx, sensor)
}

func (s *Sensor) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	sensors, err := s.sensorRepo.GetSensors(ctx)

	return sensors, err
}

func (s *Sensor) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	sensor, err := s.sensorRepo.GetSensorByID(ctx, id)

	return sensor, err
}
