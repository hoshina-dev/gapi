package repository

import (
	"context"
	"fmt"

	"github.com/hoshina-dev/gapi/internal/adapters/repository/models"
	"github.com/hoshina-dev/gapi/internal/core/domain"
	"github.com/hoshina-dev/gapi/internal/core/ports"
	"gorm.io/gorm"
)

type osmLineRepository struct {
	db *gorm.DB
}

func NewOSMLineRepository(db *gorm.DB) ports.OSMLineRepository {
	return &osmLineRepository{db: db}
}

const osmLineSearchQuery = `
WITH filtered_lines AS (
  SELECT name, way
  FROM planet_osm_line
  WHERE name IS NOT NULL
    AND tags->'name:en' IS NOT NULL
    AND (COALESCE(name,'') || ' ' || COALESCE(tags->'name:en','')) ILIKE $1
  LIMIT $2
)
SELECT name, ST_AsGeoJSON(ST_Transform(way, 4326)) AS geom
FROM filtered_lines;
`

const defaultOSMLineLimit = 20

// SearchByName implements ports.OSMLineRepository.
func (r *osmLineRepository) SearchByName(ctx context.Context, searchTerm string, limit int) ([]*domain.OSMLine, error) {
	if limit <= 0 {
		limit = defaultOSMLineLimit
	}

	return searchByName(r.db, ctx, searchTerm, limit)
}

// searchByName executes the OSM line search query and returns domain models
func searchByName(db *gorm.DB, ctx context.Context, searchTerm string, limit int) ([]*domain.OSMLine, error) {
	var queryResults []models.OSMLineSearchQuery

	// Format the search term with wildcards for ILIKE
	searchPattern := fmt.Sprintf("%%%s%%", searchTerm)

	err := db.WithContext(ctx).Raw(osmLineSearchQuery, searchPattern, limit).Scan(&queryResults).Error
	if err != nil {
		return nil, err
	}

	// Convert query results to domain models
	results := make([]*domain.OSMLine, len(queryResults))
	for i, qr := range queryResults {
		results[i] = qr.ToDomain()
	}

	return results, nil
}
