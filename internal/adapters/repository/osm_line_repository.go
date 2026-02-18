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

	// Run EXPLAIN ANALYZE to see query plan
	explainAnalyze(db, ctx, searchPattern, limit)

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

	var queryResults []models.OSMLineSearchQuery
	for rows.Next() {
		var qr models.OSMLineSearchQuery
		err := rows.Scan(&qr.Name, &qr.NameEn, &qr.Geometry, &qr.Centroid)
		if err != nil {
			return nil, err
		}
		queryResults = append(queryResults, qr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert query results to domain models
	results := make([]*domain.OSMLine, len(queryResults))
	for i, qr := range queryResults {
		results[i] = qr.ToDomain()
	}

	return results, nil
}

// explainAnalyze runs EXPLAIN ANALYZE and logs the query plan
func explainAnalyze(db *gorm.DB, ctx context.Context, searchPattern string, limit int) {
	explainQuery := fmt.Sprintf("EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) %s", osmLineSearchQuery)

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
		err := rows.Scan(&result)
		if err != nil {
			log.Printf("[EXPLAIN ANALYZE SCAN ERROR] %v", err)
			return
		}
		log.Printf("[EXPLAIN ANALYZE RESULT]\n%s\n", result)
	}
}
