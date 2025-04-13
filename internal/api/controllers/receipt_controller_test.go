package controllers

import (
	"billing/internal/core/services"
	"billing/internal/db/ent"
	"billing/internal/db/ent/enttest"
	"billing/internal/db/repositories/mocks"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTest(t *testing.T) (*ReceiptController, *mocks.MockReceiptRepository, *ent.Client) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	mockRepo := new(mocks.MockReceiptRepository)
	service := services.NewReceiptService(client, mockRepo)
	controller := NewReceiptController(service)
	return controller, mockRepo, client
}

func TestProcessReceipt(t *testing.T) {
	controller, mockRepo, client := setupTest(t)
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	t.Run("successful receipt processing", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("payload", ReceiptPayload{
			Retailer:     "Test Retailer",
			PurchaseDate: "2023-10-01",
			PurchaseTime: "10:00",
			Total:        "100.00",
			Items: []ItemPayload{
				{ShortDescription: "Item 1", Price: "50.00"},
				{ShortDescription: "Item 2", Price: "50.00"},
			},
		})

		mockRepo.On("CreateReceiptWithItems", mock.Anything, mock.AnythingOfType("*ent.Receipt"), mock.Anything).Return(&ent.Receipt{ID: 1}, nil)

		controller.ProcessReceipt(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"id":"1"}`, w.Body.String())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid payload format", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("payload", "invalid payload")

		controller.ProcessReceipt(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid payload format"}`, w.Body.String())
	})

	t.Run("invalid total format", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("payload", ReceiptPayload{
			Retailer:     "Test Retailer",
			PurchaseDate: "2023-10-01",
			PurchaseTime: "10:00",
			Total:        "invalid",
			Items: []ItemPayload{
				{ShortDescription: "Item 1", Price: "50.00"},
			},
		})

		controller.ProcessReceipt(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Please verify input. Invalid total format"}`, w.Body.String())
	})
}

func TestGetReceiptPoints(t *testing.T) {
	controller, mockRepo, client := setupTest(t)
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	t.Run("successful points retrieval", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "1"}}

		mockRepo.On("GetReceiptByID", mock.Anything, 1).Return(&ent.Receipt{ID: 1, Points: 100}, nil)

		controller.GetReceiptPoints(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"points":100}`, w.Body.String())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid receipt ID format", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		controller.GetReceiptPoints(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"Invalid receipt ID format"}`, w.Body.String())
	})

	t.Run("receipt not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{{Key: "id", Value: "999"}}

		mockRepo.On("GetReceiptByID", mock.Anything, 999).Return((*ent.Receipt)(nil), fmt.Errorf("receipt not found"))

		controller.GetReceiptPoints(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"Receipt not found"}`, w.Body.String())
	})
}
