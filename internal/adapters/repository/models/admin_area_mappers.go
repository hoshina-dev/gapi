package models

import (
	"github.com/hoshina-dev/gapi/internal/core/domain"
)

func newDomain(id int, name, iso string, level int32, parent *string, geom []byte) *domain.AdminArea {
	return &domain.AdminArea{ID: id, Name: name, ISOCode: iso, AdminLevel: level, ParentCode: parent, Geometry: geom}
}

func (a AdminArea0) ToDomain() *domain.AdminArea {
	return newDomain(a.ID, a.Name, a.GID0, 0, nil, a.Geometry)
}

func (a AdminArea1) ToDomain() *domain.AdminArea {
	return newDomain(a.ID, a.Name, a.GID1, 1, &a.GID0, a.Geometry)
}

func (a AdminArea2) ToDomain() *domain.AdminArea {
	return newDomain(a.ID, a.Name, a.GID2, 2, &a.GID1, a.Geometry)
}

func (a AdminArea3) ToDomain() *domain.AdminArea {
	return newDomain(a.ID, a.Name, a.GID3, 3, &a.GID2, a.Geometry)
}

func (a AdminArea4) ToDomain() *domain.AdminArea {
	return newDomain(a.ID, a.Name, a.GID4, 4, &a.GID3, a.Geometry)
}

func MapAdminSliceToDomain[T AdminArea](in []T) []*domain.AdminArea {
	out := make([]*domain.AdminArea, len(in))
	for i := range in {
		out[i] = in[i].ToDomain()
	}
	return out
}
