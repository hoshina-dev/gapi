package models

import (
	"github.com/hoshina-dev/gapi/internal/core/domain"
)

type OSMLineSearchQuery struct {
	Name     string `gorm:"column:name"`
	Geometry []byte `gorm:"column:geom"`
}

// ToDomain converts OSMLineSearchQuery to domain model
func (q OSMLineSearchQuery) ToDomain() *domain.OSMLine {
	return &domain.OSMLine{
		Name:     q.Name,
		Geometry: q.Geometry,
	}
}
