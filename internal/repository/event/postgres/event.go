package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"homework/internal/domain"
	"time"
)

var ErrEventNotFound = errors.New("event not found")

type EventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		pool,
	}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if event == nil {
		return fmt.Errorf("event is nil")
	}

	const sql = `
        INSERT INTO events (timestamp, sensor_serial_number, sensor_id, payload)
        VALUES ($1, $2, $3, $4)
    `

	_, err := r.pool.Exec(ctx,
		sql,
		event.Timestamp,
		event.SensorSerialNumber,
		event.SensorID,
		event.Payload,
	)

	return err
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const sqlQuery = `
        SELECT (timestamp, sensor_serial_number, sensor_id, payload)
        FROM events
        WHERE sensor_id = $1
        ORDER BY timestamp DESC
        LIMIT 1
    `

	row := r.pool.QueryRow(ctx, sqlQuery, id)

	var e domain.Event
	if err := row.Scan(&e.Timestamp, &e.SensorSerialNumber, &e.SensorID, &e.Payload); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no events yet on sensor %d: %w", id, ErrEventNotFound)
		}
		return nil, err
	}

	return &e, nil
}

func (r *EventRepository) GetEventsInRangeBySensorID(ctx context.Context, id int64, from, to time.Time) ([]*domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const sqlQuery = `
        SELECT (timestamp, sensor_serial_number, sensor_id, payload)
        FROM events
        WHERE sensor_id = $1 AND timestamp BETWEEN $2 AND $3 
        ORDER BY timestamp
    `

	rows, err := r.pool.Query(ctx, sqlQuery, id, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.Timestamp, &e.SensorSerialNumber, &e.SensorID, &e.Payload); err != nil {
			return nil, err
		}
		result = append(result, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
