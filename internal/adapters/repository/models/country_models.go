package models

import (
	"encoding/json"

	"github.com/hoshina-dev/gapi/internal/core/domain"
)

type Country struct {
	ID       int    `gorm:"column:ogc_fid;primaryKey"`
	Name     string `gorm:"column:country"`
	ISOCode  string `gorm:"column:gid_0"`
	Geometry []byte `gorm:"column:geom;type:geometry(MultiPolygon,4326)"`
}

func (c Country) ToDomain() (*domain.Country, error) {
	var geom map[string]any
	if err := json.Unmarshal(c.Geometry, &geom); err != nil {
		return nil, err
	}

	country := &domain.Country{
		ID:       c.ID,
		Name:     c.Name,
		ISOCode:  c.ISOCode,
		Geometry: geom,
	}
	return country, nil
}
