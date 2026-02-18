package domain

// OSMLine represents an OSM line feature from planet_osm_line table
type OSMLine struct {
	Name     *string     `json:"name"`
	NameEn   *string     `json:"name_en"`
	Geometry []byte      `json:"geom"`
	Centroid *Coordinate `json:"centroid"`
}
