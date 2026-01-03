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
func (c *cacheAdminAreaRepository) GetByID(ctx context.Context, id int, adminLevel int32) (*domain.AdminArea, error) {
	cacheKey := fmt.Sprintf("admin_area:%d:%d", adminLevel, id)

	var adminArea domain.AdminArea
	if c.cache.Get(ctx, cacheKey, &adminArea) {
		return &adminArea, nil
	}

	// Cache miss: fetch from underlying repo
	result, err := c.repo.GetByID(ctx, id, adminLevel)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ctx, cacheKey, result)
	return result, nil
}

// List implements ports.AdminAreaRepository.
func (c *cacheAdminAreaRepository) List(ctx context.Context, adminLevel int32) ([]*domain.AdminArea, error) {
	cacheKey := fmt.Sprintf("admin_area:list:%d", adminLevel)

	var adminAreas []*domain.AdminArea
	if c.cache.Get(ctx, cacheKey, &adminAreas) {
		return adminAreas, nil
	}

	// Cache miss: fetch from underlying repo
	result, err := c.repo.List(ctx, adminLevel)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ctx, cacheKey, result)
	return result, nil
}

// GetByCode implements ports.AdminAreaRepository.
func (c *cacheAdminAreaRepository) GetByCode(ctx context.Context, code string, adminLevel int32) (*domain.AdminArea, error) {
	cacheKey := fmt.Sprintf("admin_area:code:%d:%s", adminLevel, code)

	var adminArea domain.AdminArea
	if c.cache.Get(ctx, cacheKey, &adminArea) {
		return &adminArea, nil
	}

	// Cache miss: fetch from underlying repo
	result, err := c.repo.GetByCode(ctx, code, adminLevel)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ctx, cacheKey, result)
	return result, nil
}

// GetChildren implements ports.AdminAreaRepository.
func (c *cacheAdminAreaRepository) GetChildren(ctx context.Context, parentCode string, childLevel int32) ([]*domain.AdminArea, error) {
	cacheKey := fmt.Sprintf("admin_area:children:%d:%s", childLevel, parentCode)

	var adminAreas []*domain.AdminArea
	if c.cache.Get(ctx, cacheKey, &adminAreas) {
		return adminAreas, nil
	}

	// Cache miss: fetch from underlying repo
	result, err := c.repo.GetChildren(ctx, parentCode, childLevel)
	if err != nil {
		return nil, err
	}

	c.cache.Set(ctx, cacheKey, result)
	return result, nil
}
