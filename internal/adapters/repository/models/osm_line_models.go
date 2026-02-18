package models

import (
	"encoding/json"

	"github.com/hoshina-dev/gapi/internal/core/domain"
)

type OSMLineSearchQuery struct {
	Name     *string `gorm:"column:name"`
	NameEn   *string `gorm:"column:name_en"`
	Geometry []byte  `gorm:"column:geom"`
	Centroid []byte  `gorm:"column:centroid"`
}

// GeoJSONPoint represents a GeoJSON point geometry
type GeoJSONPoint struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

// ToDomain converts OSMLineSearchQuery to domain model
func (q OSMLineSearchQuery) ToDomain() *domain.OSMLine {
	var centroidCoord *domain.Coordinate

	if len(q.Centroid) > 0 {
		var geoJSON GeoJSONPoint
		if err := json.Unmarshal(q.Centroid, &geoJSON); err == nil {
			// GeoJSON coordinates are [lon, lat]
			centroidCoord = &domain.Coordinate{
				Lat: geoJSON.Coordinates[1],
				Lon: geoJSON.Coordinates[0],
			}
		}
	}

	return &domain.OSMLine{
		Name:     q.Name,
		NameEn:   q.NameEn,
		Geometry: q.Geometry,
		Centroid: centroidCoord,
	}
}
