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

const osmLineWithAddressQuery = `
WITH road AS (
    SELECT
        name,
        tags->'name:en' AS name_en,
        way,
        ST_Transform(way, 4326) AS geom_4326
    FROM planet_osm_line
    WHERE (name IS NOT NULL AND tags ? 'name:en')
      AND (COALESCE(name,'') || ' ' || COALESCE(tags->'name:en','')) ILIKE $1
    LIMIT $2
)
SELECT
    r.name,
    r.name_en,
    ST_AsGeoJSON(r.geom_4326) AS geom,
    ST_AsGeoJSON(ST_LineInterpolatePoint(r.geom_4326, 0.5)) AS centroid,
    COALESCE(a4.name_4, a3.name_3, a2.name_2, a1.name_1, a0.country) AS admin4,
    COALESCE(a3.name_3, a2.name_2, a1.name_1, a0.country) AS admin3,
    COALESCE(a2.name_2, a1.name_1, a0.country) AS admin2,
    COALESCE(a1.name_1, a0.country) AS admin1,
    COALESCE(a0.country, 'Unknown') AS country
FROM road r
LEFT JOIN LATERAL (SELECT name_4 FROM admin4 WHERE geom && r.geom_4326 AND ST_Intersects(geom, r.geom_4326) LIMIT 1) a4 ON TRUE
LEFT JOIN LATERAL (SELECT name_3 FROM admin3 WHERE geom && r.geom_4326 AND ST_Intersects(geom, r.geom_4326) LIMIT 1) a3 ON TRUE
LEFT JOIN LATERAL (SELECT name_2 FROM admin2 WHERE geom && r.geom_4326 AND ST_Intersects(geom, r.geom_4326) LIMIT 1) a2 ON TRUE
LEFT JOIN LATERAL (SELECT name_1 FROM admin1 WHERE geom && r.geom_4326 AND ST_Intersects(geom, r.geom_4326) LIMIT 1) a1 ON TRUE
LEFT JOIN LATERAL (SELECT country FROM admin0 WHERE geom && r.geom_4326 AND ST_Intersects(geom, r.geom_4326) LIMIT 1) a0 ON TRUE;
`

const defaultOSMLineLimit = 20

// SearchRoadName implements ports.OSMLineRepository.
func (r *osmLineRepository) SearchRoadName(ctx context.Context, searchTerm string, limit int) ([]*domain.OSMLine, error) {
	if limit <= 0 {
		limit = defaultOSMLineLimit
	}

	return searchRoadName(r.db, ctx, searchTerm, limit)
}

// GetAddressByRoadName searches for OSM lines by name and returns address information
func (r *osmLineRepository) GetAddressByRoadName(ctx context.Context, searchTerm string, limit int) ([]*domain.LineWithAddress, error) {
	if limit <= 0 {
		limit = defaultOSMLineLimit
	}

	return getAddressByRoadName(r.db, ctx, searchTerm, limit)
}

// searchRoadName executes the OSM line search query and returns domain models
func searchRoadName(db *gorm.DB, ctx context.Context, searchTerm string, limit int) ([]*domain.OSMLine, error) {
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

func getAddressByRoadName(db *gorm.DB, ctx context.Context, searchTerm string, limit int) ([]*domain.LineWithAddress, error) {
	searchPattern := fmt.Sprintf("%%%s%%", searchTerm)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	rows, err := sqlDB.QueryContext(ctx, osmLineWithAddressQuery, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.LineWithAddress
	for rows.Next() {
		var qr models.OSMLineAddressQuery
		err := rows.Scan(&qr.Name, &qr.NameEn, &qr.Geometry, &qr.Centroid, &qr.Admin4, &qr.Admin3, &qr.Admin2, &qr.Admin1, &qr.Country)
		if err != nil {
			return nil, err
		}
		results = append(results, qr.ToDomain())
	}

	return results, nil
}
