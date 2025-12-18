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
func (c *adminAreaService) GetAll(ctx context.Context, adminLevel int32) ([]*domain.AdminArea, error) {
	return c.repo.List(ctx, adminLevel)
}

// GetByID implements [ports.AdminAreaService].
func (c *adminAreaService) GetByID(ctx context.Context, id int, adminLevel int32) (*domain.AdminArea, error) {
	return c.repo.GetByID(ctx, id, adminLevel)
}

// GetByCode implements [ports.AdminAreaService].
func (c *adminAreaService) GetByCode(ctx context.Context, code string, admin_level int32) (*domain.AdminArea, error) {
	return c.repo.GetByCode(ctx, code, admin_level)
}

// GetChildren implements [ports.AdminAreaService].
func (c *adminAreaService) GetChildren(ctx context.Context, parentCode string, adminLevel int32) ([]*domain.AdminArea, error) {
	panic("unimplemented")
}
