package services

import (
	"context"

	"github.com/hoshina-dev/gapi/internal/core/domain"
	"github.com/hoshina-dev/gapi/internal/core/ports"
)

type countryService struct {
	repo ports.CountryRepository
}

func NewCountryService(repo ports.CountryRepository) ports.CountryService {
	return &countryService{repo: repo}
}

// GetAll implements [ports.CountryService].
func (c *countryService) GetAll(ctx context.Context) ([]*domain.Country, error) {
	admin_level := 0
	return c.repo.List(ctx, &admin_level)
}

// GetByID implements [ports.CountryService].
func (c *countryService) GetByID(ctx context.Context, id int) (*domain.Country, error) {
	return c.repo.GetByID(ctx, id)
}

// GetByCode implements [ports.CountryService].
func (c *countryService) GetByCode(ctx context.Context, code string) (*domain.Country, error) {
	return c.repo.GetByCode(ctx, code, 1)
}
