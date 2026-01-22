package ports

import (
	"context"

	"github.com/hoshina-dev/gapi/internal/core/domain"
)

type AdminAreaService interface {
	GetAll(ctx context.Context, adminLevel int32, tolerance *float64) ([]*domain.AdminArea, error)
	GetByID(ctx context.Context, id int, adminLevel int32, tolerance *float64) (*domain.AdminArea, error)
	GetByCode(ctx context.Context, code string, adminLevel int32, tolerance *float64) (*domain.AdminArea, error)
	GetChildren(ctx context.Context, parentCode string, childLevel int32, tolerance *float64) ([]*domain.AdminArea, error)
	FilterCoordinatesByBoundary(ctx context.Context, coordinates [][2]float64, boundaryID string, adminLevel int32) ([][]float64, error)
}
