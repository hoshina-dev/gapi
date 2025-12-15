package domain

type Country struct {
	ID       int    `json:"id" gorm:"column:ogc_fid;primaryKey"`
	Name     string `json:"name" gorm:"column:country"`
	ISOCode  string `json:"iso_code" gorm:"column:gid_0"`
	Geometry string `json:"geom" gorm:"column:geom;type:geometry(MultiPolygon,4326)"`
}
