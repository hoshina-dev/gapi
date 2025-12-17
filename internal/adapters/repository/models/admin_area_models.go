package models

import (
	"encoding/json"

	"github.com/hoshina-dev/gapi/internal/core/domain"
)

type AdminArea struct {
	ID         int    `gorm:"column:ogc_fid;primaryKey"`
	Name       string `gorm:"column:country"`
	ISOCode    string `gorm:"column:gid_0"`
	AdminLevel int32  `gorm:"column:admin_level"`
	ParentID   *int   `gorm:"column:parent_id"`
	Geometry   []byte `gorm:"column:geom;type:geometry(MultiPolygon,4326)"`
}

func (a AdminArea) ToDomain() (*domain.AdminArea, error) {
	var geom map[string]any
	if err := json.Unmarshal(a.Geometry, &geom); err != nil {
		return nil, err
	}

	adminArea := &domain.AdminArea{
		ID:         a.ID,
		Name:       a.Name,
		ISOCode:    a.ISOCode,
		Geometry:   geom,
		AdminLevel: a.AdminLevel,
		ParentID:   a.ParentID,
	}
	return adminArea, nil
}
