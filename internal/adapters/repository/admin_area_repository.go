package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

// GetByID implements ports.AdminAreaRepository.
func (c *adminAreaRepository) GetByID(ctx context.Context, id int, adminLevel int32) (*domain.AdminArea, error) {
	query, ok := queries[adminLevel]
	if !ok {
		return nil, errors.New("invalid admin level")
	}

	cacheKey := fmt.Sprintf("admin_area:%d:%d", adminLevel, id)

	// Check cache if Redis is available
	if c.redisClient != nil {
		cachedData, err := c.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var adminArea domain.AdminArea
			if json.Unmarshal([]byte(cachedData), &adminArea) == nil {
				return &adminArea, nil // Cache hit
			}
		} else if err != redis.Nil {
			// Redis error, but continue to DB
		}
	}

	// Cache miss or no Redis: fetch from DB
	var adminArea *domain.AdminArea
	switch adminLevel {
	case 0:
		var model models.AdminArea0
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).First(&model, id).Error
		if err != nil {
			return nil, err
		}
		adminArea = model.ToDomain()
	case 1:
		var model models.AdminArea1
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).First(&model, id).Error
		if err != nil {
			return nil, err
		}
		adminArea = model.ToDomain()
	default:
		return nil, errors.New("invalid admin level")
	}

	// Store in cache if Redis is available
	if c.redisClient != nil {
		if data, marshalErr := json.Marshal(adminArea); marshalErr == nil {
			c.redisClient.Set(ctx, cacheKey, data, 0) // No TTL, rely on LRU eviction
		}
	}

	return adminArea, nil
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

	switch adminLevel {
	case 0:
		var adminArea models.AdminArea0
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
			Where("gid_0 = ?", code).First(&adminArea).Error
		if err != nil {
			return nil, err
		}
		return adminArea.ToDomain(), nil
	case 1:
		var adminArea models.AdminArea1
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
			Where("gid_1 = ?", code).First(&adminArea).Error
		if err != nil {
			return nil, err
		}
		return adminArea.ToDomain(), nil
	default:
		return nil, errors.New("invalid admin level")
	}
}

// GetChildren implements [ports.AdminAreaRepository].
func (c *adminAreaRepository) GetChildren(ctx context.Context, parentCode string, childLevel int32) ([]*domain.AdminArea, error) {
	query, ok := queries[childLevel]
	if !ok {
		return nil, errors.New("invalid child level")
	}

	switch childLevel {
	case 1:
		var adminAreas []models.AdminArea1
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
			Where("gid_0 = ?", parentCode).
			Order("name_1").Scan(&adminAreas).Error
		if err != nil {
			return nil, err
		}
		return models.MapAdmin1SliceToDomain(adminAreas), nil
	default:
		return nil, errors.New("invalid child level")
	}
}

func (c *adminAreaRepository) listAdmin0(ctx context.Context) ([]*domain.AdminArea, error) {
	query := queries[0]
	var adminAreas []models.AdminArea0
	err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
		Order("country").Scan(&adminAreas).Error
	if err != nil {
		return nil, err
	}
	return models.MapAdmin0SliceToDomain(adminAreas), nil
}

func (c *adminAreaRepository) listAdmin1(ctx context.Context) ([]*domain.AdminArea, error) {
	query := queries[1]
	var adminAreas []models.AdminArea1
	err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
		Order("name_1").Scan(&adminAreas).Error
	if err != nil {
		return nil, err
	}
	return models.MapAdmin1SliceToDomain(adminAreas), nil
}
