package repository

import (
	"context"

	"github.com/hoshina-dev/gapi/internal/adapters/repository/models"
	"github.com/hoshina-dev/gapi/internal/core/domain"
	"github.com/hoshina-dev/gapi/internal/core/ports"
	"gorm.io/gorm"
)

type countryRepository struct {
	db *gorm.DB
}

func NewCountryRepository(db *gorm.DB) ports.CountryRepository {
	return &countryRepository{db: db}
}

// GetByID implements ports.CountryRepository.
func (c *countryRepository) GetByID(ctx context.Context, id int) (*domain.Country, error) {
	var country *models.AdminArea

	err := c.db.WithContext(ctx).Table("admin_areas").
		Select("ogc_fid", "gid_0", "country", "admin_level", "parent_id", "ST_AsGeoJSON(geom) AS geom").
		First(&country, id).Error
	if err != nil {
		return nil, err
	}

	return country.ToDomain()
}

// List implements ports.CountryRepository.
func (c *countryRepository) List(ctx context.Context) ([]*domain.Country, error) {
	var results []*models.AdminArea

	err := c.db.WithContext(ctx).Table("admin_areas").
		Select("ogc_fid", "gid_0", "country", "admin_level", "parent_id", "ST_AsGeoJSON(geom) AS geom").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	countries := make([]*domain.Country, len(results))
	for i, res := range results {
		country, err := res.ToDomain()
		if err != nil {
			return nil, err
		}
		countries[i] = country
	}

	return countries, err
}

// GetByCode implements [ports.CountryRepository].
func (c *countryRepository) GetByCode(ctx context.Context, code string) (*domain.Country, error) {
	var country *models.AdminArea

	err := c.db.WithContext(ctx).Table("admin_areas").
		Select("ogc_fid", "gid_0", "country", "admin_level", "parent_id", "ST_AsGeoJSON(geom) AS geom").
		Where("gid_0 = ?", code).First(&country).Error
	if err != nil {
		return nil, err
	}

	return country.ToDomain()
}
