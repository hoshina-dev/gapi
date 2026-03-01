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

type OSMLineAddressQuery struct {
	Name     *string `gorm:"column:name"`
	NameEn   *string `gorm:"column:name_en"`
	Geometry []byte  `gorm:"column:geom"`
	Centroid []byte  `gorm:"column:centroid"`
	Admin4   *string `gorm:"column:admin4"`
	Admin3   *string `gorm:"column:admin3"`
	Admin2   *string `gorm:"column:admin2"`
	Admin1   *string `gorm:"column:admin1"`
	Country  *string `gorm:"column:country"`
}

// GeoJSONPoint represents a GeoJSON point geometry
type GeoJSONPoint struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

// ToDomain converts OSMLineSearchQuery to domain model
func (q OSMLineSearchQuery) ToDomain() *domain.OSMLine {
	var centroidCoord domain.Coordinate // zero-value = empty Coordinate{} when absent

	if len(q.Centroid) > 0 {
		var geoJSON GeoJSONPoint
		if err := json.Unmarshal(q.Centroid, &geoJSON); err == nil {
			// GeoJSON coordinates are [lon, lat]
			centroidCoord = domain.Coordinate{
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

// ToDomainWithAddress converts OSMLineAddressQuery to domain model with address
func (q OSMLineAddressQuery) ToDomain() *domain.LineWithAddress {
	var centroidCoord domain.Coordinate // zero-value = empty Coordinate{} when absent

	if len(q.Centroid) > 0 {
		var geoJSON GeoJSONPoint
		if err := json.Unmarshal(q.Centroid, &geoJSON); err == nil {
			// GeoJSON coordinates are [lon, lat]
			centroidCoord = domain.Coordinate{
				Lat: geoJSON.Coordinates[1],
				Lon: geoJSON.Coordinates[0],
			}
		}
	}

	return &domain.LineWithAddress{
		Line: domain.OSMLine{
			Name:     q.Name,
			NameEn:   q.NameEn,
			Geometry: q.Geometry,
			Centroid: centroidCoord,
		},
		Address: &domain.AdminAddress{
			Country: q.Country,
			Admin1:  q.Admin1,
			Admin2:  q.Admin2,
			Admin3:  q.Admin3,
			Admin4:  q.Admin4,
		},
	}
}
