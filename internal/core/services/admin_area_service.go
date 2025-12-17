package services

import (
	"context"

	"github.com/hoshina-dev/gapi/internal/core/domain"
	"github.com/hoshina-dev/gapi/internal/core/ports"
)

type adminAreaService struct {
	repo ports.AdminAreaRepository
}

func NewAdminAreaService(repo ports.AdminAreaRepository) ports.AdminAreaService {
	return &adminAreaService{repo: repo}
}

// GetAll implements [ports.AdminAreaService].
func (c *adminAreaService) GetAll(ctx context.Context, admin_level *int32) ([]*domain.AdminArea, error) {
	return c.repo.List(ctx, admin_level)
}

// GetByID implements [ports.AdminAreaService].
func (c *adminAreaService) GetByID(ctx context.Context, id int) (*domain.AdminArea, error) {
	return c.repo.GetByID(ctx, id)
}

// GetByCode implements [ports.AdminAreaService].
func (c *adminAreaService) GetByCode(ctx context.Context, code string, admin_level int32) (*domain.AdminArea, error) {
	return c.repo.GetByCode(ctx, code, admin_level)
}
