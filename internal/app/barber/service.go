package barber

import (
	"context"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

type BarberService struct {
	barberRepo BarberRepository
}

func NewService(barberRepo BarberRepository) *BarberService {
	return &BarberService{
		barberRepo: barberRepo,
	}
}

func (s *BarberService) GetByID(ctx context.Context, id int64) (domain.Barber, error) {
	return s.barberRepo.GetByID(ctx, id)
}

func (s *BarberService) ListAll(ctx context.Context) ([]domain.Barber, error) {
	return s.barberRepo.ListAll(ctx)
}

func (s *BarberService) SetActive(ctx context.Context, id int64, active bool) error {
	return s.barberRepo.SetActive(ctx, id, active)
}
