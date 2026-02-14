package domain

type AdminArea struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	ISOCode    string  `json:"iso_code"`
	AdminLevel int32   `json:"admin_level"`
	ParentCode *string `json:"parent_code"`
	Geometry   []byte  `json:"geom"`
}

type Coordinate struct {
	ID  string
	Lat float64
	Lon float64
}

type FilteredCoordinate struct {
	Idx int
	Lat float64
	Lon float64
}
