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
func (c *adminAreaService) GetAll(ctx context.Context, adminLevel int32, tolerance *float64) ([]*domain.AdminArea, error) {
	return c.repo.List(ctx, adminLevel, tolerance)
}

// GetByID implements [ports.AdminAreaService].
func (c *adminAreaService) GetByID(ctx context.Context, id int, adminLevel int32, tolerance *float64) (*domain.AdminArea, error) {
	return c.repo.GetByID(ctx, id, adminLevel, tolerance)
}

// GetByCode implements [ports.AdminAreaService].
func (c *adminAreaService) GetByCode(ctx context.Context, code string, adminLevel int32, tolerance *float64) (*domain.AdminArea, error) {
	return c.repo.GetByCode(ctx, code, adminLevel, tolerance)
}

// GetChildren implements [ports.AdminAreaService].
func (c *adminAreaService) GetChildren(ctx context.Context, parentCode string, childLevel int32, tolerance *float64) ([]*domain.AdminArea, error) {
	return c.repo.GetChildren(ctx, parentCode, childLevel, tolerance)
}

// FilterCoordinatesByBoundary implements [ports.AdminAreaService].
func (c *adminAreaService) FilterCoordinatesByBoundary(ctx context.Context, coordinates [][2]float64, boundaryID string, adminLevel int32) ([][]float64, error) {
	return c.repo.FilterCoordinatesByBoundary(ctx, coordinates, boundaryID, adminLevel)
}
