package repository

import (
	"context"

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
	panic("unimplemented")
}

// List implements ports.CountryRepository.
func (c *countryRepository) List(ctx context.Context) ([]domain.Country, error) {
	var countries []domain.Country

	err := c.db.WithContext(ctx).Raw("SELECT ogc_fid, gid_0, country, ST_AsGeoJSON(geom) AS geom FROM countries").Scan(&countries).Error

	return countries, err
}
