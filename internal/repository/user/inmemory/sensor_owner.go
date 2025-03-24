package inmemory

import (
	"context"
	"homework/internal/domain"
	"sync"
)

type SensorOwnerRepository struct {
	mu             sync.RWMutex
	usersToSensors map[int64]map[int64]struct{}
}

func NewSensorOwnerRepository() *SensorOwnerRepository {
	return &SensorOwnerRepository{
		usersToSensors: make(map[int64]map[int64]struct{}),
	}
}

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, sensorOwner domain.SensorOwner) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	sensors, ok := r.usersToSensors[sensorOwner.UserID]
	if !ok {
		sensors = make(map[int64]struct{})
		r.usersToSensors[sensorOwner.UserID] = sensors
	}

	sensors[sensorOwner.SensorID] = struct{}{}

	return nil
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	sensors, ok := r.usersToSensors[userID]
	if !ok {
		return []domain.SensorOwner{}, nil
	}

	owners := make([]domain.SensorOwner, 0, len(sensors))
	for sensorID := range sensors {
		owners = append(owners, domain.SensorOwner{
			UserID:   userID,
			SensorID: sensorID,
		})
	}

	return owners, nil
}
