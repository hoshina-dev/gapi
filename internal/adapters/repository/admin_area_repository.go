package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hoshina-dev/gapi/internal/adapters/repository/models"
	"github.com/hoshina-dev/gapi/internal/core/domain"
	"github.com/hoshina-dev/gapi/internal/core/ports"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type adminAreaRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewAdminAreaRepository(db *gorm.DB, redisClient *redis.Client) ports.AdminAreaRepository {
	return &adminAreaRepository{db: db, redisClient: redisClient}
}

var queries = map[int32]struct{ Table, Select string }{
	0: {"admin0", "ogc_fid, gid_0, country, ST_AsGeoJSON(geom) AS geom"},
	1: {"admin1", "ogc_fid, gid_0, gid_1, name_1, ST_AsGeoJSON(geom) AS geom"},
}

// Helper methods for caching
func (c *adminAreaRepository) getFromCache(ctx context.Context, cacheKey string, dest interface{}) bool {
	if c.redisClient == nil {
		return false
	}
	data, err := c.redisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		// Only log non-cache-miss errors
		if err != redis.Nil {
			log.Warnf("Redis Get error for key %s: %v", cacheKey, err)
		}
		return false
	}
	if err := json.Unmarshal([]byte(data), dest); err != nil {
		log.Warnf("Failed to unmarshal cached data for key %s: %v", cacheKey, err)
		return false
	}
	return true
}

func (c *adminAreaRepository) setToCache(ctx context.Context, cacheKey string, value interface{}) {
	if c.redisClient == nil {
		return
	}
	data, err := json.Marshal(value)
	if err != nil {
		log.Warnf("Failed to marshal data for cache key %s: %v", cacheKey, err)
		return
	}
	if err := c.redisClient.Set(ctx, cacheKey, data, 0).Err(); err != nil {
		log.Warnf("Failed to set cache for key %s: %v", cacheKey, err)
	}
}

// GetByID implements ports.AdminAreaRepository.
func (c *adminAreaRepository) GetByID(ctx context.Context, id int, adminLevel int32) (*domain.AdminArea, error) {
	query, ok := queries[adminLevel]
	if !ok {
		return nil, errors.New("invalid admin level")
	}

	cacheKey := fmt.Sprintf("admin_area:%d:%d", adminLevel, id)

	var adminArea domain.AdminArea
	if c.getFromCache(ctx, cacheKey, &adminArea) {
		return &adminArea, nil
	}

	// Cache miss or no Redis: fetch from DB
	switch adminLevel {
	case 0:
		var model models.AdminArea0
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).First(&model, id).Error
		if err != nil {
			return nil, err
		}
		adminArea = *model.ToDomain()
	case 1:
		var model models.AdminArea1
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).First(&model, id).Error
		if err != nil {
			return nil, err
		}
		adminArea = *model.ToDomain()
	default:
		return nil, errors.New("invalid admin level")
	}

	c.setToCache(ctx, cacheKey, &adminArea)
	return &adminArea, nil
}

// List implements ports.AdminAreaRepository.
func (c *adminAreaRepository) List(ctx context.Context, adminLevel int32) ([]*domain.AdminArea, error) {
	switch adminLevel {
	case 0:
		return c.listAdmin0(ctx)
	case 1:
		return c.listAdmin1(ctx)
	default:
		return nil, errors.New("invalid admin level")
	}
}

// GetByCode implements [ports.AdminAreaRepository].
func (c *adminAreaRepository) GetByCode(ctx context.Context, code string, adminLevel int32) (*domain.AdminArea, error) {
	query, ok := queries[adminLevel]
	if !ok {
		return nil, errors.New("invalid admin level")
	}

	cacheKey := fmt.Sprintf("admin_area:code:%d:%s", adminLevel, code)

	var adminArea domain.AdminArea
	if c.getFromCache(ctx, cacheKey, &adminArea) {
		return &adminArea, nil
	}

	// Cache miss or no Redis: fetch from DB
	switch adminLevel {
	case 0:
		var model models.AdminArea0
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
			Where("gid_0 = ?", code).First(&model).Error
		if err != nil {
			return nil, err
		}
		adminArea = *model.ToDomain()
	case 1:
		var model models.AdminArea1
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
			Where("gid_1 = ?", code).First(&model).Error
		if err != nil {
			return nil, err
		}
		adminArea = *model.ToDomain()
	default:
		return nil, errors.New("invalid admin level")
	}

	c.setToCache(ctx, cacheKey, &adminArea)
	return &adminArea, nil
}

// GetChildren implements [ports.AdminAreaRepository].
func (c *adminAreaRepository) GetChildren(ctx context.Context, parentCode string, childLevel int32) ([]*domain.AdminArea, error) {
	query, ok := queries[childLevel]
	if !ok {
		return nil, errors.New("invalid child level")
	}

	cacheKey := fmt.Sprintf("admin_area:children:%d:%s", childLevel, parentCode)

	var adminAreas []*domain.AdminArea
	if c.getFromCache(ctx, cacheKey, &adminAreas) {
		return adminAreas, nil
	}

	// Cache miss or no Redis: fetch from DB
	switch childLevel {
	case 1:
		var adminModels []models.AdminArea1
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
			Where("gid_0 = ?", parentCode).
			Order("name_1").Scan(&adminModels).Error
		if err != nil {
			return nil, err
		}
		adminAreas = models.MapAdmin1SliceToDomain(adminModels)
	default:
		return nil, errors.New("invalid child level")
	}

	c.setToCache(ctx, cacheKey, adminAreas)
	return adminAreas, nil
}

func (c *adminAreaRepository) listAdmin0(ctx context.Context) ([]*domain.AdminArea, error) {
	cacheKey := "admin_area:list:0"

	var adminAreas []*domain.AdminArea
	if c.getFromCache(ctx, cacheKey, &adminAreas) {
		return adminAreas, nil
	}

	// Cache miss or no Redis: fetch from DB
	query := queries[0]
	var adminModels []models.AdminArea0
	err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
		Order("country").Scan(&adminModels).Error
	if err != nil {
		return nil, err
	}
	adminAreas = models.MapAdmin0SliceToDomain(adminModels)

	c.setToCache(ctx, cacheKey, adminAreas)
	return adminAreas, nil
}

func (c *adminAreaRepository) listAdmin1(ctx context.Context) ([]*domain.AdminArea, error) {
	cacheKey := "admin_area:list:1"

	var adminAreas []*domain.AdminArea
	if c.getFromCache(ctx, cacheKey, &adminAreas) {
		return adminAreas, nil
	}

	// Cache miss or no Redis: fetch from DB
	query := queries[1]
	var adminModels []models.AdminArea1
	err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
		Order("name_1").Scan(&adminModels).Error
	if err != nil {
		return nil, err
	}
	adminAreas = models.MapAdmin1SliceToDomain(adminModels)

	c.setToCache(ctx, cacheKey, adminAreas)
	return adminAreas, nil
}
