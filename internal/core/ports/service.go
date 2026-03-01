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
	FilterCoordinatesByBoundary(ctx context.Context, coordinates []*domain.Coordinate, boundaryID string, adminLevel int32) ([]*domain.Coordinate, error)
}

type OSMLineService interface {
	SearchRoadName(ctx context.Context, searchTerm string, limit int) ([]*domain.OSMLine, error)
	GetAddressByRoadName(ctx context.Context, searchTerm string, limit int) ([]*domain.LineWithAddress, error)
	FindNearbyRoads(ctx context.Context, lat float64, lon float64, radius float64, limit int) ([]*domain.OSMLine, error)
}
