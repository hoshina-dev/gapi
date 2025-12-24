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
		0,
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

	req := httptest.NewRequest("GET", "/query?query={adminArea(id:1, adminLevel:0){id}}", nil)

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
