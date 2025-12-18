package models

import (
	"github.com/hoshina-dev/gapi/internal/core/domain"
)

func (a AdminArea0) ToDomain() *domain.AdminArea {
	return &domain.AdminArea{
		ID:         a.ID,
		Name:       a.Name,
		ISOCode:    a.GID0,
		AdminLevel: 0,
		Geometry:   a.Geometry,
	}
}

func (a AdminArea1) ToDomain() *domain.AdminArea {
	return &domain.AdminArea{
		ID:         a.ID,
		Name:       a.Name,
		ISOCode:    a.GID1,
		AdminLevel: 1,
		ParentCode: &a.GID0,
		Geometry:   a.Geometry,
	}
}
