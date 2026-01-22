package graph

import (
	"testing"

	"github.com/hoshina-dev/gapi/internal/adapters/graph/model"
	"github.com/stretchr/testify/assert"
)

func TestParseBoundaryID(t *testing.T) {
	tests := []struct {
		name        string
		boundaryID  string
		expectError bool
		expectLevel int32
		expectGID   string
	}{
		{
			name:        "valid level 0 - Thailand",
			boundaryID:  "THA",
			expectError: false,
			expectLevel: 0,
			expectGID:   "THA",
		},
		{
			name:        "valid level 0 - Austria",
			boundaryID:  "AUT",
			expectError: false,
			expectLevel: 0,
			expectGID:   "AUT",
		},
		{
			name:        "valid level 1",
			boundaryID:  "AUT.1",
			expectError: false,
			expectLevel: 1,
			expectGID:   "AUT.1",
		},
		{
			name:        "valid level 2",
			boundaryID:  "AUT.1.4",
			expectError: false,
			expectLevel: 2,
			expectGID:   "AUT.1.4",
		},
		{
			name:        "valid level 3",
			boundaryID:  "AUT.1.4.15",
			expectError: false,
			expectLevel: 3,
			expectGID:   "AUT.1.4.15",
		},
		{
			name:        "valid level 4",
			boundaryID:  "AUT.1.4.15.1",
			expectError: false,
			expectLevel: 4,
			expectGID:   "AUT.1.4.15.1",
		},
		{
			name:        "valid with underscore in level",
			boundaryID:  "THA.3_1",
			expectError: false,
			expectLevel: 1,
			expectGID:   "THA.3_1",
		},
		{
			name:        "empty string",
			boundaryID:  "",
			expectError: true,
			expectLevel: 0,
		},
		{
			name:        "invalid ISO code - too short",
			boundaryID:  "T",
			expectError: true,
			expectLevel: 0,
		},
		{
			name:        "invalid ISO code - too long",
			boundaryID:  "THAI",
			expectError: true,
			expectLevel: 0,
		},
		{
			name:        "too many levels",
			boundaryID:  "AUT.1.2.3.4.5",
			expectError: true,
			expectLevel: 0,
		},
		{
			name:        "valid 2-char ISO code",
			boundaryID:  "TH",
			expectError: false,
			expectLevel: 0,
			expectGID:   "TH",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseBoundaryID(tt.boundaryID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectLevel, result.AdminLevel)
				assert.Equal(t, tt.expectGID, result.GIDValue)
				assert.Equal(t, tt.boundaryID, result.FullID)
			}
		})
	}
}

func TestValidateCoordinates(t *testing.T) {
	tests := []struct {
		name        string
		coordinates []*model.CoordinateInput
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid single coordinate",
			coordinates: []*model.CoordinateInput{
				{Lat: 13.7563, Lon: 100.5018},
			},
			expectError: false,
		},
		{
			name: "valid multiple coordinates",
			coordinates: []*model.CoordinateInput{
				{Lat: 13.7563, Lon: 100.5018},
				{Lat: 48.2082, Lon: 16.3738},
				{Lat: -33.8688, Lon: 151.2093},
			},
			expectError: false,
		},
		{
			name: "valid boundary values",
			coordinates: []*model.CoordinateInput{
				{Lat: 90, Lon: 180},
				{Lat: -90, Lon: -180},
				{Lat: 0, Lon: 0},
			},
			expectError: false,
		},
		{
			name:        "empty array",
			coordinates: []*model.CoordinateInput{},
			expectError: true,
			errorMsg:    "coordinates array cannot be empty",
		},
		{
			name:        "nil array",
			coordinates: nil,
			expectError: true,
			errorMsg:    "coordinates array cannot be empty",
		},
		{
			name: "latitude too high",
			coordinates: []*model.CoordinateInput{
				{Lat: 91, Lon: 100.5018},
			},
			expectError: true,
			errorMsg:    "invalid latitude at index 0: must be between -90 and 90",
		},
		{
			name: "latitude too low",
			coordinates: []*model.CoordinateInput{
				{Lat: -91, Lon: 100.5018},
			},
			expectError: true,
			errorMsg:    "invalid latitude at index 0: must be between -90 and 90",
		},
		{
			name: "longitude too high",
			coordinates: []*model.CoordinateInput{
				{Lat: 13.7563, Lon: 181},
			},
			expectError: true,
			errorMsg:    "invalid longitude at index 0: must be between -180 and 180",
		},
		{
			name: "longitude too low",
			coordinates: []*model.CoordinateInput{
				{Lat: 13.7563, Lon: -181},
			},
			expectError: true,
			errorMsg:    "invalid longitude at index 0: must be between -180 and 180",
		},
		{
			name: "nil coordinate in array",
			coordinates: []*model.CoordinateInput{
				{Lat: 13.7563, Lon: 100.5018},
				nil,
			},
			expectError: true,
			errorMsg:    "coordinate at index 1 cannot be nil",
		},
		{
			name: "invalid coordinate at second index",
			coordinates: []*model.CoordinateInput{
				{Lat: 13.7563, Lon: 100.5018},
				{Lat: 13.7563, Lon: 100.5018},
				{Lat: 95, Lon: 100.5018},
			},
			expectError: true,
			errorMsg:    "invalid latitude at index 2: must be between -90 and 90",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCoordinates(tt.coordinates)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCoordinates_ArraySizeLimit(t *testing.T) {
	t.Run("exactly at limit (10,000)", func(t *testing.T) {
		coordinates := make([]*model.CoordinateInput, 10000)
		for i := 0; i < 10000; i++ {
			coordinates[i] = &model.CoordinateInput{Lat: 0, Lon: 0}
		}

		err := validateCoordinates(coordinates)
		assert.NoError(t, err)
	})

	t.Run("exceeds limit (10,001)", func(t *testing.T) {
		coordinates := make([]*model.CoordinateInput, 10001)
		for i := 0; i < 10001; i++ {
			coordinates[i] = &model.CoordinateInput{Lat: 0, Lon: 0}
		}

		err := validateCoordinates(coordinates)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "coordinates array cannot exceed 10,000 items")
	})
}

func TestValidateTolerance(t *testing.T) {
	tests := []struct {
		name        string
		tolerance   *float64
		expectError bool
		expectNil   bool
	}{
		{
			name:        "nil tolerance",
			tolerance:   nil,
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "zero tolerance",
			tolerance:   floatPtr(0),
			expectError: false,
			expectNil:   true,
		},
		{
			name:        "positive tolerance",
			tolerance:   floatPtr(0.001),
			expectError: false,
			expectNil:   false,
		},
		{
			name:        "negative tolerance",
			tolerance:   floatPtr(-0.001),
			expectError: true,
			expectNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateTolerance(tt.tolerance)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "tolerance must be non-negative")
			} else {
				assert.NoError(t, err)
				if tt.expectNil {
					assert.Nil(t, result)
				} else {
					assert.NotNil(t, result)
					assert.Equal(t, *tt.tolerance, *result)
				}
			}
		})
	}
}

func floatPtr(f float64) *float64 {
	return &f
}
