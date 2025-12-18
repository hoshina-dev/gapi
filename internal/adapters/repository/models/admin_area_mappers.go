package models

import (
	"encoding/json"

	"github.com/hoshina-dev/gapi/internal/core/domain"
)

func (a AdminArea0) ToDomain() (*domain.AdminArea, error) {
	var geom map[string]any
	if err := json.Unmarshal(a.Geometry, &geom); err != nil {
		return nil, err
	}

	adminArea := &domain.AdminArea{
		ID:         a.ID,
		Name:       a.Name,
		ISOCode:    a.GID0,
		AdminLevel: 0,
		Geometry:   geom,
	}
	return adminArea, nil
}

func (a AdminArea1) ToDomain() (*domain.AdminArea, error) {
	var geom map[string]any
	if err := json.Unmarshal(a.Geometry, &geom); err != nil {
		return nil, err
	}

	adminArea := &domain.AdminArea{
		ID:         a.ID,
		Name:       a.Name,
		ISOCode:    a.GID1,
		AdminLevel: 1,
		ParentCode: &a.GID0,
		Geometry:   geom,
	}
	return adminArea, nil
}
