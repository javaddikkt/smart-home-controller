package postgres

import (
	"context"
	"homework/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SensorOwnerRepository struct {
	pool *pgxpool.Pool
}

func NewSensorOwnerRepository(pool *pgxpool.Pool) *SensorOwnerRepository {
	return &SensorOwnerRepository{pool}
}

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, owner domain.SensorOwner) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	const sql = `
    INSERT INTO sensors_users (user_id, sensor_id)
    VALUES ($1, $2)
  `
	_, err := r.pool.Exec(ctx, sql, owner.UserID, owner.SensorID)
	return err
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const sqlQuery = `
        SELECT sensor_id
          FROM sensors_users
         WHERE user_id = $1
         ORDER BY sensor_id
    `

	rows, _ := r.pool.Query(ctx, sqlQuery, userID)
	defer rows.Close()

	var owners []domain.SensorOwner
	for rows.Next() {
		var sensorID int64
		if err := rows.Scan(&sensorID); err != nil {
			return nil, err
		}
		owners = append(owners, domain.SensorOwner{
			UserID:   userID,
			SensorID: sensorID,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return owners, nil
}
