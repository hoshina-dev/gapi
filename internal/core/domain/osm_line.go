package domain

// OSMLine represents an OSM line feature from planet_osm_line table
type OSMLine struct {
	Name     string `json:"name"`
	Geometry []byte `json:"geom"`
}
