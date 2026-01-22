package repository

import (
	"context"
	"fmt"

	"github.com/hoshina-dev/gapi/internal/adapters/infrastructure"
	"github.com/hoshina-dev/gapi/internal/core/domain"
	"github.com/hoshina-dev/gapi/internal/core/ports"
)

type cacheAdminAreaRepository struct {
	repo  ports.AdminAreaRepository
	cache *infrastructure.Cache
}

func NewCacheAdminAreaRepository(repo ports.AdminAreaRepository, cache *infrastructure.Cache) ports.AdminAreaRepository {
	return &cacheAdminAreaRepository{repo: repo, cache: cache}
}

// GetByID implements ports.AdminAreaRepository.
func (c *cacheAdminAreaRepository) GetByID(ctx context.Context, id int, adminLevel int32, tolerance *float64) (*domain.AdminArea, error) {
	cacheKey := c.generateCacheKey("admin_area", adminLevel, id, tolerance)

	var adminArea domain.AdminArea
	if c.cache.Get(ctx, cacheKey, &adminArea) {
		return &adminArea, nil
	}

	// Cache miss: fetch from underlying repo
	result, err := c.repo.GetByID(ctx, id, adminLevel, tolerance)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ctx, cacheKey, result)
	return result, nil
}

// List implements ports.AdminAreaRepository.
func (c *cacheAdminAreaRepository) List(ctx context.Context, adminLevel int32, tolerance *float64) ([]*domain.AdminArea, error) {
	cacheKey := c.generateCacheKey("admin_area:list", adminLevel, tolerance)

	var adminAreas []*domain.AdminArea
	if c.cache.Get(ctx, cacheKey, &adminAreas) {
		return adminAreas, nil
	}

	// Cache miss: fetch from underlying repo
	result, err := c.repo.List(ctx, adminLevel, tolerance)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ctx, cacheKey, result)
	return result, nil
}

// GetByCode implements ports.AdminAreaRepository.
func (c *cacheAdminAreaRepository) GetByCode(ctx context.Context, code string, adminLevel int32, tolerance *float64) (*domain.AdminArea, error) {
	cacheKey := c.generateCacheKey("admin_area:code", adminLevel, code, tolerance)

	var adminArea domain.AdminArea
	if c.cache.Get(ctx, cacheKey, &adminArea) {
		return &adminArea, nil
	}

	// Cache miss: fetch from underlying repo
	result, err := c.repo.GetByCode(ctx, code, adminLevel, tolerance)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ctx, cacheKey, result)
	return result, nil
}

// GetChildren implements ports.AdminAreaRepository.
func (c *cacheAdminAreaRepository) GetChildren(ctx context.Context, parentCode string, childLevel int32, tolerance *float64) ([]*domain.AdminArea, error) {
	cacheKey := c.generateCacheKey("admin_area:children", childLevel, parentCode, tolerance)

	var adminAreas []*domain.AdminArea
	if c.cache.Get(ctx, cacheKey, &adminAreas) {
		return adminAreas, nil
	}

	// Cache miss: fetch from underlying repo
	result, err := c.repo.GetChildren(ctx, parentCode, childLevel, tolerance)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ctx, cacheKey, result)
	return result, nil
}

// FilterCoordinatesByBoundary implements ports.AdminAreaRepository.
// Note: We don't cache filtered coordinate results as they're too specific to be reused effectively.
func (c *cacheAdminAreaRepository) FilterCoordinatesByBoundary(ctx context.Context, coordinates [][2]float64, boundaryID string, adminLevel int32) ([][]float64, error) {
	// Pass through to underlying repository without caching
	// Coordinate filtering results are too specific to cache effectively
	return c.repo.FilterCoordinatesByBoundary(ctx, coordinates, boundaryID, adminLevel)
}

// generateCacheKey creates a consistent cache key by properly formatting the tolerance pointer
func (c *cacheAdminAreaRepository) generateCacheKey(prefix string, parts ...interface{}) string {
	key := prefix
	for _, part := range parts {
		switch v := part.(type) {
		case *float64:
			if v == nil {
				key += ":<nil>"
			} else {
				key += fmt.Sprintf(":%.10f", *v)
			}
		default:
			key += fmt.Sprintf(":%v", v)
		}
	}
	return key
}
