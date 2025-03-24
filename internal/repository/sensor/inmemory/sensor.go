package inmemory

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"sync"
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
	if sensor == nil {
		return fmt.Errorf("sensor is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.sensorsById[sensor.ID]; ok {
		return nil
	}

	r.sensorsById[sensor.ID] = sensor
	r.sensorsBySerialNumber[sensor.SerialNumber] = sensor

	return nil
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sensors := make([]domain.Sensor, 0, len(r.sensorsById))
	for _, sensor := range r.sensorsById {
		sensors = append(sensors, *sensor)
	}

	return nil, nil
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sensor, ok := r.sensorsById[id]
	if !ok {
		return nil, fmt.Errorf("sensor with id %d not found", id)
	}

	return sensor, nil
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sensor, ok := r.sensorsBySerialNumber[sn]
	if !ok {
		return nil, fmt.Errorf("sensor with serial number %s not found", sn)
	}

	return sensor, nil
}
