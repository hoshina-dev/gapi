package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

var queries = map[int32]struct{ Table, Select, OrderBy string }{
	0: {"admin0", "ogc_fid, gid_0, country, ST_AsGeoJSON(geom) AS geom", "country"},
	1: {"admin1", "ogc_fid, gid_0, gid_1, name_1, ST_AsGeoJSON(geom) AS geom", "name_1"},
	2: {"admin2", "ogc_fid, gid_0, gid_1, gid_2, name_2, ST_AsGeoJSON(geom) AS geom", "name_2"},
	3: {"admin3", "ogc_fid, gid_0, gid_1, gid_2, gid_3, name_3, ST_AsGeoJSON(geom) AS geom", "name_3"},
	4: {"admin4", "ogc_fid, gid_0, gid_1, gid_2, gid_3, gid_4, name_4, ST_AsGeoJSON(geom) AS geom", "name_4"},
}

// GetByID implements ports.AdminAreaRepository.
func (c *adminAreaRepository) GetByID(ctx context.Context, id int, adminLevel int32, tolerance *float64) (*domain.AdminArea, error) {
	switch adminLevel {
	case 0:
		return getByID[models.AdminArea0](c.db, ctx, id, adminLevel, tolerance)
	case 1:
		return getByID[models.AdminArea1](c.db, ctx, id, adminLevel, tolerance)
	case 2:
		return getByID[models.AdminArea2](c.db, ctx, id, adminLevel, tolerance)
	case 3:
		return getByID[models.AdminArea3](c.db, ctx, id, adminLevel, tolerance)
	case 4:
		return getByID[models.AdminArea4](c.db, ctx, id, adminLevel, tolerance)
	default:
		return nil, errors.New("invalid admin level")
	}
}

// List implements ports.AdminAreaRepository.
func (c *adminAreaRepository) List(ctx context.Context, adminLevel int32, tolerance *float64) ([]*domain.AdminArea, error) {
	switch adminLevel {
	case 0:
		return list[models.AdminArea0](c.db, ctx, adminLevel, tolerance)
	case 1:
		return list[models.AdminArea1](c.db, ctx, adminLevel, tolerance)
	case 2:
		return list[models.AdminArea2](c.db, ctx, adminLevel, tolerance)
	case 3:
		return list[models.AdminArea3](c.db, ctx, adminLevel, tolerance)
	case 4:
		return list[models.AdminArea4](c.db, ctx, adminLevel, tolerance)
	default:
		return nil, errors.New("invalid admin level")
	}
}

// GetByCode implements [ports.AdminAreaRepository].
func (c *adminAreaRepository) GetByCode(ctx context.Context, code string, adminLevel int32, tolerance *float64) (*domain.AdminArea, error) {
	switch adminLevel {
	case 0:
		return getByCode[models.AdminArea0](c.db, ctx, code, adminLevel, tolerance)
	case 1:
		return getByCode[models.AdminArea1](c.db, ctx, code, adminLevel, tolerance)
	case 2:
		return getByCode[models.AdminArea2](c.db, ctx, code, adminLevel, tolerance)
	case 3:
		return getByCode[models.AdminArea3](c.db, ctx, code, adminLevel, tolerance)
	case 4:
		return getByCode[models.AdminArea4](c.db, ctx, code, adminLevel, tolerance)
	default:
		return nil, errors.New("invalid admin level")
	}
}

// GetChildren implements [ports.AdminAreaRepository].
func (c *adminAreaRepository) GetChildren(ctx context.Context, parentCode string, childLevel int32, tolerance *float64) ([]*domain.AdminArea, error) {
	switch childLevel {
	case 1:
		return getChildren[models.AdminArea1](c.db, ctx, parentCode, childLevel, tolerance)
	case 2:
		return getChildren[models.AdminArea2](c.db, ctx, parentCode, childLevel, tolerance)
	case 3:
		return getChildren[models.AdminArea3](c.db, ctx, parentCode, childLevel, tolerance)
	case 4:
		return getChildren[models.AdminArea4](c.db, ctx, parentCode, childLevel, tolerance)
	default:
		return nil, errors.New("invalid child level")
	}
}

func getByID[T models.AdminArea](db *gorm.DB, ctx context.Context, id int, adminLevel int32, tolerance *float64) (*domain.AdminArea, error) {
	query := queries[adminLevel]
	var adminArea T
	selectClause := getSelectClause(query.Select, tolerance)
	q := db.WithContext(ctx).Table(query.Table).Select(selectClause)
	if err := q.First(&adminArea, id).Error; err != nil {
		return nil, err
	}
	return adminArea.ToDomain(), nil
}

func list[T models.AdminArea](db *gorm.DB, ctx context.Context, adminLevel int32, tolerance *float64) ([]*domain.AdminArea, error) {
	query := queries[adminLevel]
	var adminAreas []T
	selectClause := getSelectClause(query.Select, tolerance)
	q := db.WithContext(ctx).Table(query.Table).Select(selectClause)
	if err := q.Order(query.OrderBy).Scan(&adminAreas).Error; err != nil {
		return nil, err
	}
	return models.MapAdminSliceToDomain(adminAreas), nil
}

func getByCode[T models.AdminArea](db *gorm.DB, ctx context.Context, code string, adminLevel int32, tolerance *float64) (*domain.AdminArea, error) {
	query := queries[adminLevel]
	gidCol := "gid_" + strconv.Itoa(int(adminLevel))
	var adminArea T
	selectClause := getSelectClause(query.Select, tolerance)
	q := db.WithContext(ctx).Table(query.Table).Select(selectClause)
	if strings.Contains(code, "_") {
		if err := q.Where(gidCol+" = ?", code).First(&adminArea).Error; err != nil {
			return nil, err
		}
	} else if err := q.Where(gidCol+" LIKE ? ESCAPE '\\' ", code+"\\_%").First(&adminArea).Error; err != nil {
		return nil, err
	}

	return adminArea.ToDomain(), nil
}

func getChildren[T models.AdminArea](db *gorm.DB, ctx context.Context, parentCode string, childLevel int32, tolerance *float64) ([]*domain.AdminArea, error) {
	query := queries[childLevel]
	whereClause := "gid_" + strconv.Itoa(int(childLevel-1)) + " = ?"
	var adminAreas []T
	selectClause := getSelectClause(query.Select, tolerance)
	q := db.WithContext(ctx).Table(query.Table).Select(selectClause)
	if err := q.Where(whereClause, parentCode).Order(query.OrderBy).Scan(&adminAreas).Error; err != nil {
		return nil, err
	}
	return models.MapAdminSliceToDomain(adminAreas), nil
}

func getSelectClause(baseSelect string, tolerance *float64) string {
	if tolerance != nil && *tolerance > 0 {
		// Replace geom with simplified geometry using tolerance value
		return strings.ReplaceAll(
			baseSelect,
			"ST_AsGeoJSON(geom) AS geom",
			fmt.Sprintf("ST_AsGeoJSON(ST_SimplifyPreserveTopology(geom, %f)) AS geom", *tolerance),
		)
	}
	return baseSelect
}

// FilterCoordinatesByBoundary implements [ports.AdminAreaRepository].
func (c *adminAreaRepository) FilterCoordinatesByBoundary(ctx context.Context, coordinates [][2]float64, boundaryID string, adminLevel int32) ([][]float64, error) {
	query := queries[adminLevel]
	if query.Table == "" {
		return nil, errors.New("invalid admin level")
	}

	// Build VALUES clause for coordinates
	// Format: (idx, lat, lon), (idx, lat, lon), ...
	valuesClauses := make([]string, len(coordinates))
	for i, coord := range coordinates {
		valuesClauses[i] = fmt.Sprintf("(%d, %f, %f)", i, coord[0], coord[1])
	}
	valuesSQL := strings.Join(valuesClauses, ", ")

	// Build SQL query using CTE
	// Uses ST_Contains to filter coordinates within the boundary polygon
	// Note: ST_MakePoint takes (lon, lat) not (lat, lon)!
	sql := fmt.Sprintf(`
		WITH
			boundary AS (
				SELECT geom FROM %s WHERE gid_%d = ?
			),
			input_coords(idx, lat, lon) AS (
				VALUES %s
			)
		SELECT c.lat, c.lon
		FROM input_coords c, boundary b
		WHERE ST_Contains(
			b.geom,
			ST_SetSRID(ST_MakePoint(c.lon, c.lat), 4326)
		)
		ORDER BY c.idx
	`, query.Table, adminLevel, valuesSQL)

	type Result struct {
		Lat float64
		Lon float64
	}

	var results []Result
	if err := c.db.WithContext(ctx).Raw(sql, boundaryID).Scan(&results).Error; err != nil {
		return nil, err
	}

	// If no results found, check if boundary exists
	if len(results) == 0 {
		var count int64
		whereClause := fmt.Sprintf("gid_%d = ?", adminLevel)
		if err := c.db.WithContext(ctx).Table(query.Table).Where(whereClause, boundaryID).Count(&count).Error; err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, fmt.Errorf("boundary not found: %s", boundaryID)
		}
	}

	// Convert results to output format
	output := make([][]float64, len(results))
	for i, r := range results {
		output[i] = []float64{r.Lat, r.Lon}
	}

	return output, nil
}
