package graph

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hoshina-dev/gapi/internal/adapters/graph/model"
)

// validateTolerance ensures tolerance is not negative and returns nil if it's 0 or less
func validateTolerance(tolerance *float64) (*float64, error) {
	if tolerance == nil {
		return nil, nil
	}
	if *tolerance < 0 {
		return nil, errors.New("tolerance must be non-negative")
	}
	if *tolerance == 0 {
		return nil, nil
	}
	return tolerance, nil
}

// BoundaryInfo contains parsed boundary ID information
type BoundaryInfo struct {
	FullID     string
	AdminLevel int32
	GIDValue   string
}

// parseBoundaryID parses and validates boundary ID format
// Format: "ISO3.level1.level2.level3.level4"
// Examples: "THA" (level 0), "AUT.1" (level 1), "AUT.1.4.15.1" (level 4)
func parseBoundaryID(boundaryID string) (*BoundaryInfo, error) {
	if boundaryID == "" {
		return nil, errors.New("boundaryId cannot be empty")
	}

	parts := strings.Split(boundaryID, ".")
	if len(parts) < 1 || len(parts) > 5 {
		return nil, errors.New("invalid boundaryId format: must have 1-5 parts separated by dots")
	}

	// Validate ISO3 code (first part) - should be 2 or 3 characters
	if len(parts[0]) < 2 || len(parts[0]) > 3 {
		return nil, errors.New("invalid boundaryId format: first part must be 2-3 character ISO code")
	}

	adminLevel := int32(len(parts) - 1)

	return &BoundaryInfo{
		FullID:     boundaryID,
		AdminLevel: adminLevel,
		GIDValue:   boundaryID,
	}, nil
}

// validateCoordinates ensures coordinates array is within limits and has valid values
func validateCoordinates(coordinates []*model.CoordinateInput) error {
	if len(coordinates) == 0 {
		return errors.New("coordinates array cannot be empty")
	}

	if len(coordinates) > 10000 {
		return errors.New("coordinates array cannot exceed 10,000 items")
	}

	for i, coord := range coordinates {
		if coord == nil {
			return fmt.Errorf("coordinate at index %d cannot be nil", i)
		}
		if coord.Lat < -90 || coord.Lat > 90 {
			return fmt.Errorf("invalid latitude at index %d: must be between -90 and 90", i)
		}
		if coord.Lon < -180 || coord.Lon > 180 {
			return fmt.Errorf("invalid longitude at index %d: must be between -180 and 180", i)
		}
	}

	return nil
}
