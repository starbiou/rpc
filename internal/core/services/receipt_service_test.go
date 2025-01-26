package services

import (
	"billing/internal/db/ent"
	"billing/internal/db/ent/enttest"
	"billing/internal/db/repositories/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	_ "github.com/mattn/go-sqlite3"
)

func setupReceiptService(t *testing.T) (*ReceiptService, *mocks.MockReceiptRepository, *ent.Client) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	mockRepo := new(mocks.MockReceiptRepository)
	service := NewReceiptService(client, mockRepo)
	return service, mockRepo, client
}

func TestAddReceipt(t *testing.T) {
	service, mockRepo, client := setupReceiptService(t)
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	receipt := &ent.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        125, // $1.25 in cents
		Edges: ent.ReceiptEdges{
			Items: []*ent.Item{
				{ShortDescription: "Pepsi - 12-oz", Price: 125},
			},
		},
	}

	expectedReceipt := &ent.Receipt{
		ID:           1,
		Retailer:     receipt.Retailer,
		PurchaseDate: receipt.PurchaseDate,
		PurchaseTime: receipt.PurchaseTime,
		Total:        receipt.Total,
	}

	mockRepo.On("CreateReceiptWithItems", mock.Anything, mock.AnythingOfType("*ent.Receipt"),
		receipt.Edges.Items).Return(expectedReceipt, nil)

	id, err := service.AddReceipt(context.Background(), *receipt)
	assert.NoError(t, err)
	assert.Equal(t, "1", id)
	mockRepo.AssertExpectations(t)
}

func TestGetReceiptPoints(t *testing.T) {
	service, mockRepo, client := setupReceiptService(t)
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	expectedReceipt := &ent.Receipt{
		ID:     1,
		Points: 100,
	}

	mockRepo.On("GetReceiptByID", mock.Anything, 1).Return(expectedReceipt, nil)

	points, err := service.GetReceiptPoints(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 100, points)
	mockRepo.AssertExpectations(t)
}

func TestCalculatePoints(t *testing.T) {
	service, _, _ := setupReceiptService(t)

	receipt := &ent.Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Total:        3535,
		Edges: ent.ReceiptEdges{
			Items: []*ent.Item{
				{ShortDescription: "Mountain Dew 12PK", Price: 649},
				{ShortDescription: "Emils Cheese Pizza", Price: 1225},
				{ShortDescription: "Knorr Creamy Chicken", Price: 126},
				{ShortDescription: "Doritos Nacho Cheese", Price: 335},
				{ShortDescription: "Klarbrunn 12-PK 12 FL OZ", Price: 1200},
			},
		},
	}

	points := service.calculatePoints(receipt)

	assert.Equal(t, 28, points)
}
