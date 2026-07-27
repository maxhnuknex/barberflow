package postgres

import (
	"context"
	"fmt"

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
