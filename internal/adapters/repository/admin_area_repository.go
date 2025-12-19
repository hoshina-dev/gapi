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

	switch adminLevel {
	case 0:
		var adminArea models.AdminArea0
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).First(&adminArea, id).Error
		if err != nil {
			return nil, err
		}
		return adminArea.ToDomain(), nil
	case 1:
		var adminArea models.AdminArea1
		err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).First(&adminArea, id).Error
		if err != nil {
			return nil, err
		}
		return adminArea.ToDomain(), nil
	default:
		return nil, errors.New("invalid admin level")
	}
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

		results := make([]*domain.AdminArea, len(adminAreas))
		for i, adminArea := range adminAreas {
			results[i] = adminArea.ToDomain()
		}

		return results, err
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

	results := make([]*domain.AdminArea, len(adminAreas))
	for i, adminArea := range adminAreas {
		results[i] = adminArea.ToDomain()
	}

	return results, err
}

func (c *adminAreaRepository) listAdmin1(ctx context.Context) ([]*domain.AdminArea, error) {
	query := queries[1]
	var adminAreas []models.AdminArea1
	err := c.db.WithContext(ctx).Table(query.Table).Select(query.Select).
		Order("name_1").Scan(&adminAreas).Error
	if err != nil {
		return nil, err
	}

	results := make([]*domain.AdminArea, len(adminAreas))
	for i, adminArea := range adminAreas {
		results[i] = adminArea.ToDomain()
	}

	return results, err
}
