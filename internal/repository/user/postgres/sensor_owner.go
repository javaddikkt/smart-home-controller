package postgres

import (
	"context"
	"fmt"
	"homework/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorOwnerRepository struct {
	pool *pgxpool.Pool
}

func NewSensorOwnerRepository(pool *pgxpool.Pool) *SensorOwnerRepository {
	return &SensorOwnerRepository{pool}
}

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, owner *domain.SensorOwner) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if owner == nil {
		return fmt.Errorf("sensor owner is nil")
	}

	const sql = `
    INSERT INTO sensor_owners (user_id, sensor_id)
    VALUES ($1, $2)
    ON CONFLICT (user_id, sensor_id) DO NOTHING
  `
	_, err := r.pool.Exec(ctx, sql, owner.UserID, owner.SensorID)
	return err
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.Sensor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const sql = `
    SELECT s.id, s.serial_number, s.type, s.current_state,
           s.description, s.is_active, s.registered_at, s.last_activity
      FROM sensors s
      JOIN sensor_owners so ON so.sensor_id = s.id
     WHERE so.user_id = $1
     ORDER BY s.id
  `
	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Sensor
	for rows.Next() {
		var s domain.Sensor
		if err := rows.Scan(
			&s.ID,
			&s.SerialNumber,
			&s.Type,
			&s.CurrentState,
			&s.Description,
			&s.IsActive,
			&s.RegisteredAt,
			&s.LastActivity,
		); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
