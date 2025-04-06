package usecase

import (
	"context"
	"homework/internal/domain"
)

type User struct {
	userRepo   UserRepository
	ownerRepo  SensorOwnerRepository
	sensorRepo SensorRepository
}

func NewUser(ur UserRepository, sor SensorOwnerRepository, sr SensorRepository) *User {
	return &User{
		userRepo:   ur,
		ownerRepo:  sor,
		sensorRepo: sr,
	}
}

func (u *User) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if u.userRepo == nil {
		return nil, ErrInvalidUserName
	}
	return user, u.userRepo.SaveUser(ctx, user)
}

func (u *User) AttachSensorToUser(ctx context.Context, userID, sensorID int64) error {
	_, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = u.sensorRepo.GetSensorByID(ctx, sensorID)
	if err != nil {
		return err
	}

	return u.ownerRepo.SaveSensorOwner(ctx, domain.SensorOwner{
		UserID:   userID,
		SensorID: sensorID,
	})
}

func (u *User) GetUserSensors(ctx context.Context, userID int64) ([]domain.Sensor, error) {
	_, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	owners, err := u.ownerRepo.GetSensorsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sensors := make([]domain.Sensor, 0, len(owners))
	for _, owner := range owners {
		sensor, err := u.sensorRepo.GetSensorByID(ctx, owner.SensorID)
		if err != nil {
			return nil, err
		}
		sensors = append(sensors, *sensor)
	}
	return sensors, nil
}
