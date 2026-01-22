package http_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hoshina-dev/gapi/internal/adapters/graph"
	"github.com/hoshina-dev/gapi/internal/adapters/graph/mocks"
	"github.com/hoshina-dev/gapi/internal/adapters/http"
	"github.com/hoshina-dev/gapi/internal/adapters/infrastructure"
	"github.com/hoshina-dev/gapi/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestApp() (*fiber.App, *mocks.MockAdminAreaService) {
	cfg := infrastructure.LoadConfig()
	mockAdminAreaService := new(mocks.MockAdminAreaService)
	resolver := graph.NewResolver(mockAdminAreaService)
	app := http.SetupRouter(resolver, cfg)
	return app, mockAdminAreaService
}

func TestGraphQLEndpoint_ValidQuery(t *testing.T) {
	// Arrange
	app, mockService := setupTestApp()

	expectedAdminArea := &domain.AdminArea{
		ID:         1,
		Name:       "Thailand",
		ISOCode:    "THA",
		AdminLevel: 0,
		Geometry:   []byte("geometry"),
	}

	mockService.On("GetByID",
		mock.Anything,
		1,
		int32(0),
	).Return(expectedAdminArea, nil)

	query := `{
        "query": "query { adminArea(id: 1, adminLevel: 0) { id name isoCode adminLevel } }"
    }`

	req := httptest.NewRequest("POST", "/query", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req, -1) // -1 disables timeout for test

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]any
	json.Unmarshal(body, &result)

	// Check GraphQL response structure
	assert.NotNil(t, result["data"])
	data := result["data"].(map[string]any)
	adminArea := data["adminArea"].(map[string]any)

	assert.EqualValues(t, 1, adminArea["id"])
	assert.Equal(t, "Thailand", adminArea["name"])
	assert.Equal(t, "THA", adminArea["isoCode"])
}

func TestGraphQLEndpoint_QueryWithVariables(t *testing.T) {
	// Arrange
	app, mockService := setupTestApp()

	expectedAdminArea := &domain.AdminArea{
		ID:         1,
		Name:       "Thailand",
		ISOCode:    "THA",
		AdminLevel: 0,
		Geometry:   []byte("geometry"),
	}

	mockService.On("GetByCode",
		mock.Anything,
		"THA",
		int32(0),
	).Return(expectedAdminArea, nil)

	query := `{
        "query": "query($code: String!, $level: Int!) { adminAreaByCode(code: $code, adminLevel: $level) { id name isoCode } }",
        "variables": {
            "code": "THA",
            "level": 0
        }
    }`

	req := httptest.NewRequest("POST", "/query", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestGraphQLEndpoint_InvalidQuery(t *testing.T) {
	// Arrange
	app, _ := setupTestApp()

	query := `{
        "query": "query { invalidField { id } }"
    }`

	req := httptest.NewRequest("POST", "/query", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]any
	json.Unmarshal(body, &result)

	// Should return errors
	assert.NotNil(t, result["errors"])
}

func TestGraphQLEndpoint_MalformedJSON(t *testing.T) {
	// Arrange
	app, _ := setupTestApp()

	malformedQuery := `{ "query": invalid json }`

	req := httptest.NewRequest("POST", "/query", strings.NewReader(malformedQuery))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]any
	json.Unmarshal(body, &result)

	// Should contain error
	assert.NotNil(t, result["errors"])
}

func TestGraphQLEndpoint_GETNotAllowed(t *testing.T) {
	// Arrange
	app, _ := setupTestApp()

	req := httptest.NewRequest("GET", "/query?query=%7BadminArea%28id%3A1%2C%20adminLevel%3A0%29%7Bid%7D%7D", nil)

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)

	// GraphQL endpoint should handle GET requests too (for introspection)
	// Status depends on whether your handler supports GET
	assert.NotEqual(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestGraphQLEndpoint_ComplexQuery(t *testing.T) {
	// Arrange
	app, mockService := setupTestApp()

	expectedAdminAreas := []*domain.AdminArea{
		{
			ID:         3,
			Name:       "BangkokMetropolis",
			ISOCode:    "THA.3_1",
			AdminLevel: 1,
			Geometry:   []byte("Bangkok"),
		},
		{
			ID:         10,
			Name:       "Chiang Mai",
			ISOCode:    "THA.10_1",
			AdminLevel: 1,
			Geometry:   []byte("Chiang Mai"),
		},
	}

	mockService.On("GetChildren",
		mock.Anything,
		"THA",
		int32(1),
	).Return(expectedAdminAreas, nil)

	query := `{
        "query": "query { childrenByCode(parentCode: \"THA\", childLevel: 1 ) { id name isoCode adminLevel } }"
    }`

	req := httptest.NewRequest("POST", "/query", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]any
	json.Unmarshal(body, &result)

	data := result["data"].(map[string]any)
	adminAreas := data["childrenByCode"].([]any)

	assert.Len(t, adminAreas, 2)

	firstAdminArea := adminAreas[0].(map[string]any)
	assert.Equal(t, "BangkokMetropolis", firstAdminArea["name"])
}

func TestGraphQLEndpoint_FilterCoordinatesByBoundary_ValidRequest(t *testing.T) {
	// Arrange
	app, mockService := setupTestApp()

	expectedCoordinates := [][]float64{
		{13.7563, 100.5018}, // Bangkok
		{18.7883, 98.9853},  // Chiang Mai
	}

	mockService.On("FilterCoordinatesByBoundary",
		mock.Anything,
		[][2]float64{
			{13.7563, 100.5018},
			{48.2082, 16.3738},
			{18.7883, 98.9853},
		},
		"THA",
		int32(0),
	).Return(expectedCoordinates, nil)

	query := `{
        "query": "query { filterCoordinatesByBoundary(coordinates: [{lat: 13.7563, lon: 100.5018}, {lat: 48.2082, lon: 16.3738}, {lat: 18.7883, lon: 98.9853}], boundaryId: \"THA\") }"
    }`

	req := httptest.NewRequest("POST", "/query", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]any
	json.Unmarshal(body, &result)

	// Check GraphQL response structure
	assert.NotNil(t, result["data"])
	data := result["data"].(map[string]any)
	coordinates := data["filterCoordinatesByBoundary"].([]any)

	assert.Len(t, coordinates, 2)
	firstCoord := coordinates[0].([]any)
	assert.Equal(t, 13.7563, firstCoord[0])
	assert.Equal(t, 100.5018, firstCoord[1])
}

func TestGraphQLEndpoint_FilterCoordinatesByBoundary_WithVariables(t *testing.T) {
	// Arrange
	app, mockService := setupTestApp()

	expectedCoordinates := [][]float64{
		{13.7563, 100.5018},
	}

	mockService.On("FilterCoordinatesByBoundary",
		mock.Anything,
		[][2]float64{
			{13.7563, 100.5018},
			{48.2082, 16.3738},
		},
		"THA",
		int32(0),
	).Return(expectedCoordinates, nil)

	query := `{
        "query": "query($coords: [CoordinateInput!]!, $boundaryId: String!) { filterCoordinatesByBoundary(coordinates: $coords, boundaryId: $boundaryId) }",
        "variables": {
            "coords": [
                {"lat": 13.7563, "lon": 100.5018},
                {"lat": 48.2082, "lon": 16.3738}
            ],
            "boundaryId": "THA"
        }
    }`

	req := httptest.NewRequest("POST", "/query", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]any
	json.Unmarshal(body, &result)

	assert.NotNil(t, result["data"])
}

func TestGraphQLEndpoint_FilterCoordinatesByBoundary_EmptyCoordinates(t *testing.T) {
	// Arrange
	app, _ := setupTestApp()

	query := `{
        "query": "query { filterCoordinatesByBoundary(coordinates: [], boundaryId: \"THA\") }"
    }`

	req := httptest.NewRequest("POST", "/query", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]any
	json.Unmarshal(body, &result)

	// Should return errors
	assert.NotNil(t, result["errors"])
	errors := result["errors"].([]any)
	firstError := errors[0].(map[string]any)
	assert.Contains(t, firstError["message"], "coordinates array cannot be empty")
}

func TestGraphQLEndpoint_FilterCoordinatesByBoundary_InvalidBoundaryID(t *testing.T) {
	// Arrange
	app, _ := setupTestApp()

	query := `{
        "query": "query { filterCoordinatesByBoundary(coordinates: [{lat: 13.7563, lon: 100.5018}], boundaryId: \"X\") }"
    }`

	req := httptest.NewRequest("POST", "/query", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]any
	json.Unmarshal(body, &result)

	// Should return errors
	assert.NotNil(t, result["errors"])
	errors := result["errors"].([]any)
	firstError := errors[0].(map[string]any)
	assert.Contains(t, firstError["message"], "invalid boundaryId format")
}

func TestGraphQLEndpoint_FilterCoordinatesByBoundary_InvalidLatitude(t *testing.T) {
	// Arrange
	app, _ := setupTestApp()

	query := `{
        "query": "query { filterCoordinatesByBoundary(coordinates: [{lat: 91, lon: 100.5018}], boundaryId: \"THA\") }"
    }`

	req := httptest.NewRequest("POST", "/query", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req, -1)

	// Assert
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]any
	json.Unmarshal(body, &result)

	// Should return errors
	assert.NotNil(t, result["errors"])
	errors := result["errors"].([]any)
	firstError := errors[0].(map[string]any)
	assert.Contains(t, firstError["message"], "invalid latitude")
}

func TestGraphQLEndpoint_FilterCoordinatesByBoundary_DifferentAdminLevels(t *testing.T) {
	tests := []struct {
		name       string
		boundaryID string
		adminLevel int32
	}{
		{"level 0 - country", "THA", 0},
		{"level 1 - province", "THA.3_1", 1},
		{"level 2 - district", "AUT.1.4", 2},
		{"level 3 - sub-district", "AUT.1.4.15", 3},
		{"level 4 - municipality", "AUT.1.4.15.1", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockService := setupTestApp()

			expectedCoordinates := [][]float64{
				{13.7563, 100.5018},
			}

			mockService.On("FilterCoordinatesByBoundary",
				mock.Anything,
				[][2]float64{{13.7563, 100.5018}},
				tt.boundaryID,
				tt.adminLevel,
			).Return(expectedCoordinates, nil)

			query := `{
                "query": "query { filterCoordinatesByBoundary(coordinates: [{lat: 13.7563, lon: 100.5018}], boundaryId: \"` + tt.boundaryID + `\") }"
            }`

			req := httptest.NewRequest("POST", "/query", strings.NewReader(query))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)

			assert.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})
	}
}
