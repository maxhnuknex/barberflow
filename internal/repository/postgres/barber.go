package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

type BarberRepository struct {
	db *pgxpool.Pool
}

func NewBarberRepository(db *pgxpool.Pool) *BarberRepository {
	return &BarberRepository{
		db: db,
	}
}

func (r *BarberRepository) GetByID(ctx context.Context, id int64) (domain.Barber, error) {
	const query = `
		SELECT id, name, active, created_at
		FROM barbers
		WHERE id = $1
	`

	var barber domain.Barber
	if err := r.db.QueryRow(ctx, query, id).Scan(
		&barber.ID,
		&barber.Name,
		&barber.Active,
		&barber.CreatedAt,
	); err != nil {
		return domain.Barber{}, fmt.Errorf("get barber by id: %w", err)
	}

	return barber, nil
}

func (r *BarberRepository) ListAll(ctx context.Context) ([]domain.Barber, error) {
	const query = `
		SELECT id, name, active, created_at
		FROM barbers
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list all barbers: %w", err)
	}
	defer rows.Close()

	barbers := make([]domain.Barber, 0)
	for rows.Next() {
		var barber domain.Barber

		if err := rows.Scan(
			&barber.ID,
			&barber.Name,
			&barber.Active,
			&barber.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan barber: %w", err)
		}

		barbers = append(barbers, barber)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate barbers: %w", err)
	}

	return barbers, nil
}

func (r *BarberRepository) GetTimeActive(
	ctx context.Context,
	id int64,
	selectedDate time.Time,
) (domain.TimeInterval, error) {
	const query = `
		SELECT
			$2::date + start_time AS starts_at,
			$2::date + end_time AS ends_at
		FROM barber_working_hours
		WHERE barber_id = $1
			AND weekday = EXTRACT(ISODOW FROM $2::date)
		ORDER BY start_time
		LIMIT 1
	`

	var timeInterval domain.TimeInterval
	if err := r.db.QueryRow(ctx, query, id, selectedDate.Format("2006-01-02")).Scan(
		&timeInterval.StartsAt,
		&timeInterval.EndsAt,
	); err != nil {
		return domain.TimeInterval{}, fmt.Errorf("get active barber time: %w", err)
	}

	return timeInterval, nil
}

func (r *BarberRepository) ListBarberByService(
	ctx context.Context,
	serviceID int64,
) ([]domain.Barber, error) {
	const query = `
		SELECT 
			b.id,
			b.name,
			b.active,
			b.created_at
		FROM barbers b
		JOIN barber_services bs
    		ON bs.barber_id = b.id
		WHERE bs.service_id = $1 AND b.active = TRUE
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query, serviceID)
	if err != nil {
		return nil, fmt.Errorf("query barbers by service: %w", err)
	}
	defer rows.Close()

	barbers := make([]domain.Barber, 0)

	for rows.Next() {
		var barber domain.Barber

		if err := rows.Scan(
			&barber.ID,
			&barber.Name,
			&barber.Active,
			&barber.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan barber: %w", err)
		}

		barbers = append(barbers, barber)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate barbers: %w", err)
	}

	return barbers, nil
}

func (r *BarberRepository) SetActive(ctx context.Context, id int64, active bool) error {
	const query = `
		UPDATE barbers
		SET active = $2
		WHERE id = $1
	`

	if _, err := r.db.Exec(ctx, query, id, active); err != nil {
		return fmt.Errorf("set barber active: %w", err)
	}

	return nil
}
