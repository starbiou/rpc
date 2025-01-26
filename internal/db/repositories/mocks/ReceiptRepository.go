package mocks

import (
	"billing/internal/db/ent"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockReceiptRepository struct {
	mock.Mock
}

func (m *MockReceiptRepository) CreateReceipt(ctx context.Context, receipt *ent.Receipt) (*ent.Receipt, error) {
	args := m.Called(ctx, receipt)
	return args.Get(0).(*ent.Receipt), args.Error(1)
}

func (m *MockReceiptRepository) CreateReceiptWithItems(ctx context.Context, receipt *ent.Receipt, items []*ent.Item) (*ent.Receipt, error) {
	args := m.Called(ctx, receipt, items)
	return args.Get(0).(*ent.Receipt), args.Error(1)
}

func (m *MockReceiptRepository) GetReceiptByID(ctx context.Context, id int) (*ent.Receipt, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*ent.Receipt), args.Error(1)
}

func (m *MockReceiptRepository) UpdateReceipt(ctx context.Context, receipt *ent.Receipt) (*ent.Receipt, error) {
	args := m.Called(ctx, receipt)
	return args.Get(0).(*ent.Receipt), args.Error(1)
}

func (m *MockReceiptRepository) DeleteReceipt(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockReceiptRepository) BeginTx(ctx context.Context) (*ent.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(*ent.Tx), args.Error(1)
}

func (m *MockReceiptRepository) CommitTx(tx *ent.Tx) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockReceiptRepository) RollbackTx(tx *ent.Tx) error {
	args := m.Called(tx)
	return args.Error(0)
}
