package middleware

import (
	"billing/internal/api/controllers"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func createTestContext(payload interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	jsonBytes, _ := json.Marshal(payload)
	context.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBytes))
	context.Request.Header.Set("Content-Type", "application/json")

	return context, recorder
}

func TestValidationMiddleware(t *testing.T) {
	validate := validator.New()

	// Base valid payload
	basePayload := controllers.ReceiptPayload{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        "125.00",
		Items: []controllers.ItemPayload{
			{ShortDescription: "Item 1", Price: "50.00"},
		},
	}

	// Use reflection to generate test cases for missing fields
	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		expectedErrors int
	}{
		{
			name:           "Valid Payload",
			payload:        basePayload,
			expectedStatus: http.StatusOK,
			expectedErrors: 0,
		},
	}

	// Reflect over the basePayload to create tests for missing fields
	val := reflect.ValueOf(basePayload)
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		testPayload := basePayload

		// Set the field to its zero value
		reflect.ValueOf(&testPayload).Elem().FieldByName(fieldName).Set(reflect.Zero(val.Field(i).Type()))

		tests = append(tests, struct {
			name           string
			payload        interface{}
			expectedStatus int
			expectedErrors int
		}{
			name:           "Missing " + fieldName,
			payload:        testPayload,
			expectedStatus: http.StatusBadRequest,
			expectedErrors: 1,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, recorder := createTestContext(tt.payload)
			handler := ValidationMiddleware[controllers.ReceiptPayload](validate)
			handler(context)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			if tt.expectedStatus == http.StatusBadRequest {
				var response ValidationResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedErrors, len(response.Errors))
			}
		})
	}
}
