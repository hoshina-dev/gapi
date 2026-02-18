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
  SELECT name, tags->'name:en' AS name_en, way
  FROM planet_osm_line
  WHERE (name IS NOT NULL AND tags ? 'name:en')
    AND (COALESCE(name,'') || ' ' || COALESCE(tags->'name:en','')) ILIKE $1
  LIMIT $2
)
SELECT name, name_en, ST_AsGeoJSON(ST_Transform(way, 4326)) AS geom, ST_AsGeoJSON(ST_Transform(ST_LineInterpolatePoint(way, 0.5), 4326)) AS centroid
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
	// Format the search term with wildcards for ILIKE
	searchPattern := fmt.Sprintf("%%%s%%", searchTerm)

	// Use db.DB() to bypass GORM parameter parsing and preserve ? operator
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	rows, err := sqlDB.QueryContext(ctx, osmLineSearchQuery, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var results []*domain.OSMLine
	for rows.Next() {
		var qr models.OSMLineSearchQuery
		err := rows.Scan(&qr.Name, &qr.NameEn, &qr.Geometry, &qr.Centroid)
		if err != nil {
			return nil, err
		}
		results = append(results, qr.ToDomain())
	}

	return results, nil
}
