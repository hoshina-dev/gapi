package repository

import (
	"context"
	"errors"

	"github.com/hoshina-dev/gapi/internal/adapters/repository/models"
	"github.com/hoshina-dev/gapi/internal/core/domain"
	"github.com/hoshina-dev/gapi/internal/core/ports"
	"gorm.io/gorm"
)

type adminAreaRepository struct {
	db *gorm.DB
}

func NewAdminAreaRepository(db *gorm.DB) ports.AdminAreaRepository {
	return &adminAreaRepository{db: db}
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

	var adminArea domain.AdminArea

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

	var adminArea domain.AdminArea

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

	return &adminArea, nil
}

// GetChildren implements [ports.AdminAreaRepository].
func (c *adminAreaRepository) GetChildren(ctx context.Context, parentCode string, childLevel int32) ([]*domain.AdminArea, error) {
	query, ok := queries[childLevel]
	if !ok {
		return nil, errors.New("invalid child level")
	}

	var adminAreas []*domain.AdminArea

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

	return adminAreas, nil
}

func (c *adminAreaRepository) listAdmin0(ctx context.Context) ([]*domain.AdminArea, error) {
	var adminAreas []*domain.AdminArea

	// Cache miss or no Redis: fetch from DB
	query := queries[0]
	var adminModels []models.AdminArea0
	err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
		Order("country").Scan(&adminModels).Error
	if err != nil {
		return nil, err
	}
	adminAreas = models.MapAdmin0SliceToDomain(adminModels)

	return adminAreas, nil
}

func (c *adminAreaRepository) listAdmin1(ctx context.Context) ([]*domain.AdminArea, error) {
	var adminAreas []*domain.AdminArea

	// Cache miss or no Redis: fetch from DB
	query := queries[1]
	var adminModels []models.AdminArea1
	err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
		Order("name_1").Scan(&adminModels).Error
	if err != nil {
		return nil, err
	}
	adminAreas = models.MapAdmin1SliceToDomain(adminModels)

	return adminAreas, nil
}
