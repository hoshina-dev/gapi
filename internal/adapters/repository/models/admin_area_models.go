package models

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
