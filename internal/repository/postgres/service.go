package postgres

import (
	"context"
	"fmt"

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

func (repo *ServiceRepository) GetByID(ctx context.Context, id int64) (domain.Service, error) {
	const query = `
		SELECT id, name, duration_minutes, price_minor_units, active, created_at
		FROM services
		WHERE id = $1
	`

	var service domain.Service
	if err := repo.db.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		&service.DurationMinutes,
		&service.PriceMinorUnits,
		&service.Active,
		&service.CreatedAt,
	); err != nil {
		return domain.Service{}, fmt.Errorf("get service by id: %w", err)
	}

	return service, nil
}

func (repo *ServiceRepository) ListAll(ctx context.Context) ([]domain.Service, error) {
	const query = `
		SELECT id, name, duration_minutes, price_minor_units, active, created_at
		FROM services
		ORDER BY id
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list all services: %w", err)
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
			&service.Active,
			&service.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan service: %w", err)
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate services: %w", err)
	}

	return services, nil
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

func (repo *ServiceRepository) SetActive(ctx context.Context, id int64, active bool) error {
	const query = `
		UPDATE services
		SET active = $2
		WHERE id = $1
	`

	if _, err := repo.db.Exec(ctx, query, id, active); err != nil {
		return fmt.Errorf("set service active: %w", err)
	}

	return nil
}
