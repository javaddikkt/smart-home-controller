package inmemory

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
	"time"
)

type SensorRepository struct {
	mu                    sync.RWMutex
	sensorsById           map[int64]*domain.Sensor
	sensorsBySerialNumber map[string]*domain.Sensor
}

func NewSensorRepository() *SensorRepository {
	return &SensorRepository{
		sensorsById:           make(map[int64]*domain.Sensor),
		sensorsBySerialNumber: make(map[string]*domain.Sensor),
	}
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if sensor == nil {
		return fmt.Errorf("sensor is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if len(sensor.SerialNumber) != 10 {
		return fmt.Errorf("sensor serial number is invalid")
	}
	if sensor.Type != domain.SensorTypeContactClosure && sensor.Type != domain.SensorTypeADC {
		return fmt.Errorf("sensor type is invalid")
	}

	if _, ok := r.sensorsBySerialNumber[sensor.SerialNumber]; ok {
		return nil
	}
	sensor.RegisteredAt = time.Now()
	r.sensorsById[sensor.ID] = sensor
	r.sensorsBySerialNumber[sensor.SerialNumber] = sensor

	return nil
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	sensors := make([]domain.Sensor, 0, len(r.sensorsById))
	for _, sensor := range r.sensorsBySerialNumber {
		sensors = append(sensors, *sensor)
	}

	return sensors, nil
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	sensor, ok := r.sensorsById[id]
	if !ok {
		return nil, fmt.Errorf("sensor with id %d not found: %w", id, usecase.ErrSensorNotFound)
	}

	return sensor, nil
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	sensor, ok := r.sensorsBySerialNumber[sn]
	if !ok {
		return nil, fmt.Errorf("sensor with serial number %s not found: %w", sn, usecase.ErrSensorNotFound)
	}

	return sensor, nil
}
