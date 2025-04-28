package postgres

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorRepository struct {
	pool *pgxpool.Pool
}

func NewSensorRepository(pool *pgxpool.Pool) *SensorRepository {
	return &SensorRepository{pool}
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if sensor == nil {
		return fmt.Errorf("sensor is nil")
	}

	const sql = `
    INSERT INTO sensors (
      id, serial_number, type,
      current_state, description,
      is_active, registered_at, last_activity
    ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
    ON CONFLICT (serial_number) DO NOTHING
  `
	_, err := r.pool.Exec(ctx,
		sql,
		sensor.ID,
		sensor.SerialNumber,
		sensor.Type,
		sensor.CurrentState,
		sensor.Description,
		sensor.IsActive,
		sensor.RegisteredAt,
		sensor.LastActivity,
	)
	return err
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const sql = `
    SELECT id, serial_number, type, current_state,
           description, is_active, registered_at, last_activity
      FROM sensors
     ORDER BY id
  `
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []domain.Sensor
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
		sensors = append(sensors, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sensors, nil
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const sql = `
    SELECT id, serial_number, type, current_state,
           description, is_active, registered_at, last_activity
      FROM sensors
     WHERE id = $1
  `
	row := r.pool.QueryRow(ctx, sql, id)

	var s domain.Sensor
	if err := row.Scan(
		&s.ID,
		&s.SerialNumber,
		&s.Type,
		&s.CurrentState,
		&s.Description,
		&s.IsActive,
		&s.RegisteredAt,
		&s.LastActivity,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("sensor %d not found: %w", id, usecase.ErrSensorNotFound)
		}
		return nil, err
	}
	return &s, nil
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const sql = `
    SELECT id, serial_number, type, current_state,
           description, is_active, registered_at, last_activity
      FROM sensors
     WHERE serial_number = $1
  `
	row := r.pool.QueryRow(ctx, sql, sn)

	var s domain.Sensor
	if err := row.Scan(
		&s.ID,
		&s.SerialNumber,
		&s.Type,
		&s.CurrentState,
		&s.Description,
		&s.IsActive,
		&s.RegisteredAt,
		&s.LastActivity,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("sensor %s not found: %w", sn, usecase.ErrSensorNotFound)
		}
		return nil, err
	}
	return &s, nil
}
