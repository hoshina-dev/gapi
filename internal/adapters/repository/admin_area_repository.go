package repository

import (
	"context"

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

// GetByID implements ports.AdminAreaRepository.
func (c *adminAreaRepository) GetByID(ctx context.Context, id int) (*domain.AdminArea, error) {
	var adminArea *models.AdminArea

	err := c.db.WithContext(ctx).Table("admin_areas").
		Select("ogc_fid", "gid_0", "country", "admin_level", "parent_id", "ST_AsGeoJSON(geom) AS geom").
		First(&adminArea, id).Error
	if err != nil {
		return nil, err
	}

	return adminArea.ToDomain()
}

// List implements ports.AdminAreaRepository.
func (c *adminAreaRepository) List(ctx context.Context, admin_level *int) ([]*domain.AdminArea, error) {
	var results []*models.AdminArea
	query := c.db.WithContext(ctx).Table("admin_areas").
		Select("ogc_fid", "gid_0", "country", "admin_level", "parent_id", "ST_AsGeoJSON(geom) AS geom")

	if admin_level != nil {
		query = query.Where("admin_level = ?", *admin_level)
	}

	err := query.Order("country").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	countries := make([]*domain.AdminArea, len(results))
	for i, res := range results {
		adminArea, err := res.ToDomain()
		if err != nil {
			return nil, err
		}
		countries[i] = adminArea
	}

	return countries, err
}

// GetByCode implements [ports.AdminAreaRepository].
func (c *adminAreaRepository) GetByCode(ctx context.Context, code string, admin_level int) (*domain.AdminArea, error) {
	var adminArea *models.AdminArea

	err := c.db.WithContext(ctx).Table("admin_areas").
		Select("ogc_fid", "gid_0", "country", "admin_level", "parent_id", "ST_AsGeoJSON(geom) AS geom").
		Where("gid_0 = ? AND admin_level = ?", code, admin_level).First(&adminArea).Error
	if err != nil {
		return nil, err
	}

	return adminArea.ToDomain()
}
