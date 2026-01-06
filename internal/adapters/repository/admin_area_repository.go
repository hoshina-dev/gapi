package repository

import (
	"context"
	"errors"
	"strconv"

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
func (c *adminAreaRepository) GetByID(ctx context.Context, id int, adminLevel int32) (*domain.AdminArea, error) {
	switch adminLevel {
	case 0:
		return getByID[models.AdminArea0](c.db, ctx, id, adminLevel)
	case 1:
		return getByID[models.AdminArea1](c.db, ctx, id, adminLevel)
	case 2:
		return getByID[models.AdminArea2](c.db, ctx, id, adminLevel)
	case 3:
		return getByID[models.AdminArea3](c.db, ctx, id, adminLevel)
	case 4:
		return getByID[models.AdminArea4](c.db, ctx, id, adminLevel)
	default:
		return nil, errors.New("invalid admin level")
	}
}

// List implements ports.AdminAreaRepository.
func (c *adminAreaRepository) List(ctx context.Context, adminLevel int32) ([]*domain.AdminArea, error) {
	switch adminLevel {
	case 0:
		return list[models.AdminArea0](c.db, ctx, adminLevel)
	case 1:
		return list[models.AdminArea1](c.db, ctx, adminLevel)
	case 2:
		return list[models.AdminArea2](c.db, ctx, adminLevel)
	case 3:
		return list[models.AdminArea3](c.db, ctx, adminLevel)
	case 4:
		return list[models.AdminArea4](c.db, ctx, adminLevel)
	default:
		return nil, errors.New("invalid admin level")
	}
}

// GetByCode implements [ports.AdminAreaRepository].
func (c *adminAreaRepository) GetByCode(ctx context.Context, code string, adminLevel int32) (*domain.AdminArea, error) {
	switch adminLevel {
	case 0:
		return getByCode[models.AdminArea0](c.db, ctx, code, adminLevel)
	case 1:
		return getByCode[models.AdminArea1](c.db, ctx, code, adminLevel)
	case 2:
		return getByCode[models.AdminArea2](c.db, ctx, code, adminLevel)
	case 3:
		return getByCode[models.AdminArea3](c.db, ctx, code, adminLevel)
	case 4:
		return getByCode[models.AdminArea4](c.db, ctx, code, adminLevel)
	default:
		return nil, errors.New("invalid admin level")
	}
}

// GetChildren implements [ports.AdminAreaRepository].
func (c *adminAreaRepository) GetChildren(ctx context.Context, parentCode string, childLevel int32) ([]*domain.AdminArea, error) {
	switch childLevel {
	case 1:
		return getChildren[models.AdminArea1](c.db, ctx, parentCode, childLevel)
	case 2:
		return getChildren[models.AdminArea2](c.db, ctx, parentCode, childLevel)
	case 3:
		return getChildren[models.AdminArea3](c.db, ctx, parentCode, childLevel)
	case 4:
		return getChildren[models.AdminArea4](c.db, ctx, parentCode, childLevel)
	default:
		return nil, errors.New("invalid child level")
	}
}

func getByID[T models.AdminArea](db *gorm.DB, ctx context.Context, id int, adminLevel int32) (*domain.AdminArea, error) {
	query := queries[adminLevel]
	var adminArea T
	err := db.WithContext(ctx).Table(query.Table).Select(query.Select).First(&adminArea, id).Error
	if err != nil {
		return nil, err
	}
	return adminArea.ToDomain(), nil
}

func list[T models.AdminArea](db *gorm.DB, ctx context.Context, adminLevel int32) ([]*domain.AdminArea, error) {
	query := queries[adminLevel]
	var adminAreas []T
	err := db.WithContext(ctx).Table(query.Table).Select(query.Select).
		Order(query.OrderBy).Scan(&adminAreas).Error
	if err != nil {
		return nil, err
	}
	return models.MapAdminSliceToDomain(adminAreas), nil
}

func getByCode[T models.AdminArea](db *gorm.DB, ctx context.Context, code string, adminLevel int32) (*domain.AdminArea, error) {
	query := queries[adminLevel]
	whereClause := "gid_" + strconv.Itoa(int(adminLevel)) + " = ?"
	var adminArea T
	err := db.WithContext(ctx).Table(query.Table).Select(query.Select).
		Where(whereClause, code).First(&adminArea).Error
	if err != nil {
		return nil, err
	}
	return adminArea.ToDomain(), nil
}

func getChildren[T models.AdminArea](db *gorm.DB, ctx context.Context, parentCode string, childLevel int32) ([]*domain.AdminArea, error) {
	query := queries[childLevel]
	whereClause := "gid_" + strconv.Itoa(int(childLevel-1)) + " = ?"
	var adminAreas []T
	err := db.WithContext(ctx).Table(query.Table).Select(query.Select).
		Where(whereClause, parentCode).
		Order(query.OrderBy).Scan(&adminAreas).Error
	if err != nil {
		return nil, err
	}
	return models.MapAdminSliceToDomain(adminAreas), nil
}
