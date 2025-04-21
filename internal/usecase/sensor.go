package usecase

import (
	"context"
	"errors"
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
	if len(sensor.SerialNumber) != 10 {
		return nil, ErrWrongSensorSerialNumber
	}
	if sensor.Type != domain.SensorTypeContactClosure && sensor.Type != domain.SensorTypeADC {
		return nil, ErrWrongSensorType
	}

	currSensor, err := s.sensorRepo.GetSensorBySerialNumber(ctx, sensor.SerialNumber)
	if err == nil || !errors.Is(err, ErrSensorNotFound) {
		return currSensor, err
	}

	return sensor, s.sensorRepo.SaveSensor(ctx, sensor)
}

func (s *Sensor) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	sensors, err := s.sensorRepo.GetSensors(ctx)

	return sensors, err
}

func (s *Sensor) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	sensor, err := s.sensorRepo.GetSensorByID(ctx, id)

	return sensor, err
}
