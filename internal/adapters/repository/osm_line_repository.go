package repository

import (
	"context"
	"fmt"
	"log"

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
  WHERE (name IS NOT NULL OR tags ? 'name:en')
		AND (COALESCE(name,'') || ' ' || COALESCE(tags->'name:en','')) ILIKE $1 ESCAPE '\'
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
    WHERE (name IS NOT NULL OR tags ? 'name:en')
	AND (COALESCE(name,'') || ' ' || COALESCE(tags->'name:en','')) ILIKE $1 ESCAPE '\'
    LIMIT $2
)
SELECT
    r.name,
    r.name_en,
    ST_AsGeoJSON(r.geom_4326) AS geom,
    ST_AsGeoJSON(ST_LineInterpolatePoint(r.geom_4326, 0.5)) AS centroid,
    a.name_4 AS admin4,
    a.name_3 AS admin3,
    a.name_2 AS admin2,
    a.name_1 AS admin1,
    a.country
FROM road r
LEFT JOIN LATERAL (
    SELECT 4 AS lvl, name_4, name_3, name_2, name_1, country FROM admin4 WHERE geom && r.geom_4326 AND ST_Intersects(geom, r.geom_4326)
    UNION ALL
    SELECT 3, NULL, name_3, name_2, name_1, country FROM admin3 WHERE geom && r.geom_4326 AND ST_Intersects(geom, r.geom_4326)
    UNION ALL
    SELECT 2, NULL, NULL, name_2, name_1, country FROM admin2 WHERE geom && r.geom_4326 AND ST_Intersects(geom, r.geom_4326)
    UNION ALL
    SELECT 1, NULL, NULL, NULL, name_1, country FROM admin1 WHERE geom && r.geom_4326 AND ST_Intersects(geom, r.geom_4326)
    UNION ALL
    SELECT 0, NULL, NULL, NULL, NULL, country FROM admin0 WHERE geom && r.geom_4326 AND ST_Intersects(geom, r.geom_4326)
    ORDER BY lvl DESC
    LIMIT 1
) a ON TRUE;
`

const osmLineNearbyQuery = `
WITH pt AS (
    SELECT ST_Transform(
        ST_SetSRID(ST_MakePoint($1, $2), 4326),
        3857
    ) AS geom
)
SELECT 
    l.name,
    l.tags->'name:en' AS name_en,
    ST_AsGeoJSON(ST_Transform(l.way, 4326)) AS geom,
    ST_AsGeoJSON(ST_Transform(ST_LineInterpolatePoint(l.way, 0.5), 4326)) AS centroid
FROM planet_osm_line l
CROSS JOIN pt
WHERE ST_DWithin(l.way, pt.geom, $3)
AND (l.name IS NOT NULL OR l.tags->'name:en' IS NOT NULL)
ORDER BY ST_Distance(l.way, pt.geom) ASC
LIMIT $4;
`

// ===== DEBUG ONLY — remove explainAnalyze and its two call sites when done =====
func explainAnalyze(db *gorm.DB, ctx context.Context, baseQuery string, searchPattern string, limit int) {
	explainQuery := fmt.Sprintf("EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) %s", baseQuery)

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("[EXPLAIN ANALYZE ERROR] %v", err)
		return
	}

	rows, err := sqlDB.QueryContext(ctx, explainQuery, searchPattern, limit)
	if err != nil {
		log.Printf("[EXPLAIN ANALYZE ERROR] %v", err)
		return
	}
	defer rows.Close()

	var result string
	if rows.Next() {
		if err := rows.Scan(&result); err != nil {
			log.Printf("[EXPLAIN ANALYZE SCAN ERROR] %v", err)
			return
		}
		log.Printf("[EXPLAIN ANALYZE RESULT]\n%s\n", result)
	}
}

// ===== END DEBUG =====

// SearchRoadName implements ports.OSMLineRepository.
func (r *osmLineRepository) SearchRoadName(ctx context.Context, searchTerm string, limit int) ([]*domain.OSMLine, error) {
	return searchRoadName(r.db, ctx, searchTerm, limit)
}

// GetAddressByRoadName searches for OSM lines by name and returns address information
func (r *osmLineRepository) GetAddressByRoadName(ctx context.Context, searchTerm string, limit int) ([]*domain.LineWithAddress, error) {
	return getAddressByRoadName(r.db, ctx, searchTerm, limit)
}

func (r *osmLineRepository) FindNearbyRoads(ctx context.Context, lat float64, lon float64, radius float64, limit int) ([]*domain.OSMLine, error) {
	return findNearbyRoads(r.db, ctx, lat, lon, radius, limit)
}

// searchRoadName executes the OSM line search query and returns domain models
func searchRoadName(db *gorm.DB, ctx context.Context, searchTerm string, limit int) ([]*domain.OSMLine, error) {
	searchPattern := fmt.Sprintf("%%%s%%", escapeLike(searchTerm))
	explainAnalyze(db, ctx, osmLineSearchQuery, searchPattern, limit) // DEBUG

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

	results := make([]*domain.OSMLine, 0, limit)
	for rows.Next() {
		var qr models.OSMLineSearchQuery
		if err := rows.Scan(&qr.Name, &qr.NameEn, &qr.Geometry, &qr.Centroid); err != nil {
			return nil, err
		}
		results = append(results, qr.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func getAddressByRoadName(db *gorm.DB, ctx context.Context, searchTerm string, limit int) ([]*domain.LineWithAddress, error) {
	searchPattern := fmt.Sprintf("%%%s%%", escapeLike(searchTerm))
	explainAnalyze(db, ctx, osmLineWithAddressQuery, searchPattern, limit) // DEBUG
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
		if err := rows.Scan(&qr.Name, &qr.NameEn, &qr.Geometry, &qr.Centroid, &qr.Admin4, &qr.Admin3, &qr.Admin2, &qr.Admin1, &qr.Country); err != nil {
			return nil, err
		}
		results = append(results, qr.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func findNearbyRoads(db *gorm.DB, ctx context.Context, lat float64, lon float64, radius float64, limit int) ([]*domain.OSMLine, error) {
	var results []*models.OSMLineSearchQuery

	if err := db.WithContext(ctx).
		Raw(osmLineNearbyQuery, lon, lat, radius, limit).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	domainResults := make([]*domain.OSMLine, len(results))
	for i, result := range results {
		domainResults[i] = result.ToDomain()
	}

	return domainResults, nil
}
