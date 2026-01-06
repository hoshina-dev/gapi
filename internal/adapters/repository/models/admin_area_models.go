package models

import "github.com/hoshina-dev/gapi/internal/core/domain"

type AdminArea interface {
	ToDomain() *domain.AdminArea
}

type AdminArea0 struct {
	ID       int    `gorm:"column:ogc_fid;primaryKey"`
	GID0     string `gorm:"column:gid_0"`
	Name     string `gorm:"column:country"`
	Geometry []byte `gorm:"column:geom;type:geometry(MultiPolygon,4326)"`
}

type AdminArea1 struct {
	ID       int    `gorm:"column:ogc_fid;primaryKey"`
	GID0     string `gorm:"column:gid_0"`
	GID1     string `gorm:"column:gid_1"`
	Name     string `gorm:"column:name_1"`
	Geometry []byte `gorm:"column:geom;type:geometry(MultiPolygon,4326)"`
}

type AdminArea2 struct {
	ID       int    `gorm:"column:ogc_fid;primaryKey"`
	GID0     string `gorm:"column:gid_0"`
	GID1     string `gorm:"column:gid_1"`
	GID2     string `gorm:"column:gid_2"`
	Name     string `gorm:"column:name_2"`
	Geometry []byte `gorm:"column:geom;type:geometry(MultiPolygon,4326)"`
}

type AdminArea3 struct {
	ID       int    `gorm:"column:ogc_fid;primaryKey"`
	GID0     string `gorm:"column:gid_0"`
	GID1     string `gorm:"column:gid_1"`
	GID2     string `gorm:"column:gid_2"`
	GID3     string `gorm:"column:gid_3"`
	Name     string `gorm:"column:name_3"`
	Geometry []byte `gorm:"column:geom;type:geometry(MultiPolygon,4326)"`
}

type AdminArea4 struct {
	ID       int    `gorm:"column:ogc_fid;primaryKey"`
	GID0     string `gorm:"column:gid_0"`
	GID1     string `gorm:"column:gid_1"`
	GID2     string `gorm:"column:gid_2"`
	GID3     string `gorm:"column:gid_3"`
	GID4     string `gorm:"column:gid_4"`
	Name     string `gorm:"column:name_4"`
	Geometry []byte `gorm:"column:geom;type:geometry(MultiPolygon,4326)"`
}
