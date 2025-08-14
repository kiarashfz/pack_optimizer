package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"pack_optimizer/internal/handler/packhandler"
	"pack_optimizer/internal/repository/sqlrepo"
	"pack_optimizer/internal/usecase/packusecase"
	"testing"

	"pack_optimizer/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestCalculatePackApi is a comprehensive integration test for the /api/v1/packs/calculate endpoint.
// It uses the real handler and use case, and repository with an in-memory database.
func TestCalculatePackApi(t *testing.T) {
	// 1. Setup the in-memory SQLite database
	gormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	// 2. Add some test packs to the in-memory database
	packs := []domain.Pack{
		{Size: 250},
		{Size: 500},
		{Size: 1000},
		{Size: 2000},
		{Size: 5000},
	}
	err = gormDB.AutoMigrate(&domain.Pack{})
	assert.NoError(t, err)
	gormDB.Create(&packs)

	// 3. Initialize the real application components with the mock database
	packRepo := sqlrepo.NewPackRepo(gormDB)
	packUseCase := packusecase.NewPackUseCase(packRepo)
	packHandler := packhandler.NewPackHandler(packUseCase)

	// 4. Setup the Fiber app with the real handler
	app := fiber.New()
	app.Post("/api/v1/packs/calculate", packHandler.CalculatePacks)

	// 5. Define a table of test cases
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "Success_ValidQuantity_501",
			requestBody:    map[string]interface{}{"quantity": 501},
			expectedStatus: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"total_items":     float64(750),
				"remaining_items": float64(249),
				"total_packs":     float64(2),
				"packs": []interface{}{
					map[string]interface{}{"size": float64(250), "count": float64(1)},
					map[string]interface{}{"size": float64(500), "count": float64(1)},
				},
			},
		},
		{
			name:           "Success_ValidQuantity_1234",
			requestBody:    map[string]interface{}{"quantity": 1234},
			expectedStatus: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"total_items":     float64(1250),
				"remaining_items": float64(16),
				"total_packs":     float64(2),
				"packs": []interface{}{
					map[string]interface{}{"size": float64(250), "count": float64(1)},
					map[string]interface{}{"size": float64(1000), "count": float64(1)},
				},
			},
		},
		{
			name:           "BadRequest_ZeroQuantity",
			requestBody:    map[string]interface{}{"quantity": 0},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "Key: 'CalculatePacksReq.Quantity' Error:Field validation for 'Quantity' failed on the 'required' tag"},
		},
		{
			name:           "BadRequest_InvalidInput_Negative",
			requestBody:    map[string]interface{}{"quantity": -10},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "Key: 'CalculatePacksReq.Quantity' Error:Field validation for 'Quantity' failed on the 'min' tag"},
		},
		{
			name:           "BadRequest_InvalidInput_String",
			requestBody:    map[string]interface{}{"quantity": "not a number"},
			expectedStatus: fiber.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "invalid request"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request body
			body, mErr := json.Marshal(tt.requestBody)
			assert.NoError(t, mErr)

			// Create a new HTTP request
			req := httptest.NewRequest("POST", "/api/v1/packs/calculate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Perform the request and get the response
			resp, testErr := app.Test(req, -1) // -1 for no timeout
			assert.NoError(t, testErr)

			// Assert the status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Assert the response body
			var responseBody map[string]interface{}
			decodeErr := json.NewDecoder(resp.Body).Decode(&responseBody)
			assert.NoError(t, decodeErr)
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
