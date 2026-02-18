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

// SearchByName implements ports.OSMLineService.
func (s *osmLineService) SearchByName(ctx context.Context, searchTerm string, limit int) ([]*domain.OSMLine, error) {
	return s.repo.SearchByName(ctx, searchTerm, limit)
}
