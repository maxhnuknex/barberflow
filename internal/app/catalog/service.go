package catalog

import (
	"context"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

type CatalogService struct {
	serviceRepo ServiceRepository
}

func NewService(serviceRepo ServiceRepository) *CatalogService {
	return &CatalogService{
		serviceRepo: serviceRepo,
	}
}

func (s *CatalogService) GetByID(ctx context.Context, id int64) (domain.Service, error) {
	return s.serviceRepo.GetByID(ctx, id)
}

func (s *CatalogService) ListAll(ctx context.Context) ([]domain.Service, error) {
	return s.serviceRepo.ListAll(ctx)
}

func (s *CatalogService) SetActive(ctx context.Context, id int64, active bool) error {
	return s.serviceRepo.SetActive(ctx, id, active)
}
