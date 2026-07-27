package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (repo *ServiceRepository) ListActive(ctx context.Context) ([]domain.Service, error) {
	const query = `
		SELECT id, name, duration_minutes, price_minor_units
		FROM services
		WHERE active = TRUE
		ORDER BY id
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := make([]domain.Service, 0)
	for rows.Next() {
		var service domain.Service

		if err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.DurationMinutes,
			&service.PriceMinorUnits,
		); err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return services, nil
}
