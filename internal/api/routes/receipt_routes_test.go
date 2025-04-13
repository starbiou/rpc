package routes

import (
	"billing/internal/db/ent"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"billing/internal/db/ent/enttest"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"

	"context"
	"fmt"
)

// setupTestEnvironment initializes the test environment and returns the router and client.
func setupTestEnvironment(t *testing.T) (*gin.Engine, *ent.Client) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	router := gin.Default()
	validatorClient := validator.New()

	config := RouteConfig{
		Router:    router,
		Client:    client,
		Validator: validatorClient,
	}

	SetupReceiptRoutes(config)
	return router, client
}

func TestProcessReceiptWithExample(t *testing.T) {
	router, client := setupTestEnvironment(t)
	defer client.Close()

	// Load example data
	data, err := os.ReadFile("../../../examples/simple-receipt.json")
	if err != nil {
		t.Fatalf("failed to read example file: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("failed to unmarshal example data: %v", err)
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestProcessReceiptWithMissingFields(t *testing.T) {
	router, client := setupTestEnvironment(t)
	defer client.Close()

	payload := map[string]interface{}{
		"retailer": "Test Retailer",
		// Missing purchaseDate and purchaseTime
		"total": "10.00",
		"items": []map[string]string{
			{"shortDescription": "Item 1", "price": "5.00"},
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProcessReceiptWithInvalidDataTypes(t *testing.T) {
	router, client := setupTestEnvironment(t)
	defer client.Close()

	payload := map[string]interface{}{
		"retailer":     "Test Retailer",
		"purchaseDate": "2023-10-10",
		"purchaseTime": "15:00",
		"total":        10.00, // Should be a string
		"items": []map[string]interface{}{
			{"shortDescription": "Item 1", "price": 5.00}, // Price should be a string
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProcessReceiptWithLargeDataSet(t *testing.T) {
	router, client := setupTestEnvironment(t)
	defer client.Close()

	items := make([]map[string]string, 1000)
	for i := 0; i < 1000; i++ {
		items[i] = map[string]string{"shortDescription": fmt.Sprintf("Item %d", i+1), "price": "1.00"}
	}

	payload := map[string]interface{}{
		"retailer":     "Test Retailer",
		"purchaseDate": "2023-10-10",
		"purchaseTime": "15:00",
		"total":        "1000.00",
		"items":        items,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetReceiptPointsForNonExistentID(t *testing.T) {
	router, client := setupTestEnvironment(t)
	defer client.Close()

	req, _ := http.NewRequest("GET", "/receipts/9999/points", nil) // Assuming 9999 does not exist
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetReceiptPoints(t *testing.T) {
	router, client := setupTestEnvironment(t)
	defer client.Close()

	// Create a receipt in the database
	receipt := client.Receipt.Create().
		SetRetailer("Test Retailer").
		SetPurchaseDate("2023-10-10").
		SetPurchaseTime("15:00").
		SetTotal(10.00).
		SaveX(context.Background())

	req, _ := http.NewRequest("GET", fmt.Sprintf("/receipts/%d/points", receipt.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
