package services

import (
	"context"

	"github.com/hoshina-dev/gapi/internal/core/domain"
	"github.com/hoshina-dev/gapi/internal/core/ports"
)

type osmLineService struct {
	repo ports.OSMLineRepository
}

func NewOSMLineService(repo ports.OSMLineRepository) ports.OSMLineService {
	return &osmLineService{repo: repo}
}

// SearchRoadName implements ports.OSMLineService.
func (s *osmLineService) SearchRoadName(ctx context.Context, searchTerm string, limit int) ([]*domain.OSMLine, error) {
	return s.repo.SearchRoadName(ctx, searchTerm, limit)
}

// GetAddressByRoadName implements ports.OSMLineService.
func (s *osmLineService) GetAddressByRoadName(ctx context.Context, searchTerm string, limit int) ([]*domain.LineWithAddress, error) {
	return s.repo.GetAddressByRoadName(ctx, searchTerm, limit)
}

// FindNearbyRoads implements ports.OSMLineService.
func (s *osmLineService) FindNearbyRoads(ctx context.Context, lat float64, lon float64, radius float64, limit int) ([]*domain.OSMLine, error) {
	return s.repo.FindNearbyRoads(ctx, lat, lon, radius, limit)
}
