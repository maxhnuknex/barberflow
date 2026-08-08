package catalog

import (
	"context"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

type ServiceRepository interface {
	GetByID(ctx context.Context, id int64) (domain.Service, error)
	ListAll(ctx context.Context) ([]domain.Service, error)
	SetActive(ctx context.Context, id int64, active bool) error
}
