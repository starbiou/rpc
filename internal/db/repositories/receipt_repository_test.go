package repositories

import (
	"billing/internal/db/ent"
	"billing/internal/db/repositories/mocks"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateReceipt(t *testing.T) {
	mockRepo := new(mocks.MockReceiptRepository)

	receipt := &ent.Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2023-10-10",
		PurchaseTime: "15:00",
		Total:        1000,
	}

	mockRepo.On("CreateReceipt", mock.Anything, receipt).Return(receipt, nil)

	createdReceipt, err := mockRepo.CreateReceipt(context.Background(), receipt)
	assert.NoError(t, err)
	assert.Equal(t, receipt.Retailer, createdReceipt.Retailer)
	mockRepo.AssertExpectations(t)
}

func TestGetReceiptByID(t *testing.T) {
	mockRepo := new(mocks.MockReceiptRepository)

	expectedReceipt := &ent.Receipt{
		ID:           1,
		Retailer:     "Test Retailer",
		PurchaseDate: "2023-10-10",
		PurchaseTime: "15:00",
		Total:        1000,
	}

	mockRepo.On("GetReceiptByID", mock.Anything, 1).Return(expectedReceipt, nil)

	receipt, err := mockRepo.GetReceiptByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, expectedReceipt.ID, receipt.ID)
	mockRepo.AssertExpectations(t)
}

func TestUpdateReceipt(t *testing.T) {
	mockRepo := new(mocks.MockReceiptRepository)

	receipt := &ent.Receipt{
		ID:           1,
		Retailer:     "Updated Retailer",
		PurchaseDate: "2023-10-11",
		PurchaseTime: "16:00",
		Total:        1500,
	}

	mockRepo.On("UpdateReceipt", mock.Anything, receipt).Return(receipt, nil)

	updatedReceipt, err := mockRepo.UpdateReceipt(context.Background(), receipt)
	assert.NoError(t, err)
	assert.Equal(t, receipt.Retailer, updatedReceipt.Retailer)
	mockRepo.AssertExpectations(t)
}

func TestDeleteReceipt(t *testing.T) {
	mockRepo := new(mocks.MockReceiptRepository)

	mockRepo.On("DeleteReceipt", mock.Anything, 1).Return(nil)

	err := mockRepo.DeleteReceipt(context.Background(), 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateReceiptWithItems(t *testing.T) {
	mockRepo := new(mocks.MockReceiptRepository)

	receipt := &ent.Receipt{
		Retailer:     "Test Retailer",
		PurchaseDate: "2023-10-10",
		PurchaseTime: "15:00",
		Total:        1000,
	}

	items := []*ent.Item{
		{
			ShortDescription: "Item 1",
			Price:            500,
		},
		{
			ShortDescription: "Item 2",
			Price:            500,
		},
	}

	mockRepo.On("CreateReceiptWithItems", mock.Anything, receipt, items).Return(receipt, nil)

	createdReceipt, err := mockRepo.CreateReceiptWithItems(context.Background(), receipt, items)
	assert.NoError(t, err)
	assert.Equal(t, receipt.Retailer, createdReceipt.Retailer)
	mockRepo.AssertExpectations(t)
}
