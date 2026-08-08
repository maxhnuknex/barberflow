package barber

import (
	"context"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

type BarberRepository interface {
	GetByID(ctx context.Context, id int64) (domain.Barber, error)
	ListAll(ctx context.Context) ([]domain.Barber, error)
	SetActive(ctx context.Context, id int64, active bool) error
}
