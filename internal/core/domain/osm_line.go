package domain

// AdminAddress represents hierarchical administrative address levels
type AdminAddress struct {
	Country *string `json:"country"` // admin 0
	Admin1  *string `json:"admin1"`  // province/state
	Admin2  *string `json:"admin2"`  // district
	Admin3  *string `json:"admin3"`  // subdistrict
	Admin4  *string `json:"admin4"`  // ward
}

// OSMLine represents an OSM line feature from planet_osm_line table
type OSMLine struct {
	Name     *string     `json:"name"`
	NameEn   *string     `json:"name_en"`
	Geometry []byte      `json:"geom"`
	Centroid *Coordinate `json:"centroid"`
}

// LineWithAddress is a composite type combining road data with administrative address information
type LineWithAddress struct {
	Line    *OSMLine      `json:"line"`
	Address *AdminAddress `json:"address"`
}
