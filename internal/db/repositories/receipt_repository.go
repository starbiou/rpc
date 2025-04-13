package repositories

import (
	"billing/internal/db/ent"
	"context"
)

// ReceiptRepository defines the methods for interacting with receipts.
type ReceiptRepository interface {
	CreateReceipt(ctx context.Context, receipt *ent.Receipt) (*ent.Receipt, error)
	GetReceiptByID(ctx context.Context, id int) (*ent.Receipt, error)
	UpdateReceipt(ctx context.Context, receipt *ent.Receipt) (*ent.Receipt, error)
	DeleteReceipt(ctx context.Context, id int) error
}

// ReceiptRepositoryImpl is the concrete implementation of ReceiptRepository.
type ReceiptRepositoryImpl struct {
	client *ent.Client
}

// NewReceiptRepository creates a new instance of ReceiptRepositoryImpl.
func NewReceiptRepository(client *ent.Client) ReceiptRepository {
	return &ReceiptRepositoryImpl{client: client}
}

func (r *ReceiptRepositoryImpl) CreateReceipt(ctx context.Context, receipt *ent.Receipt) (*ent.Receipt,
	error) {
	return r.client.Receipt.Create().
		SetRetailer(receipt.Retailer).
		SetPurchaseDate(receipt.PurchaseDate).
		SetPurchaseTime(receipt.PurchaseTime).
		SetTotal(receipt.Total).
		Save(ctx)
}

func (r *ReceiptRepositoryImpl) GetReceiptByID(ctx context.Context, id int) (*ent.Receipt, error) {
	return r.client.Receipt.Get(ctx, id)
}

func (r *ReceiptRepositoryImpl) UpdateReceipt(ctx context.Context, receipt *ent.Receipt) (*ent.Receipt,
	error) {
	return r.client.Receipt.UpdateOneID(receipt.ID).
		SetRetailer(receipt.Retailer).
		SetPurchaseDate(receipt.PurchaseDate).
		SetPurchaseTime(receipt.PurchaseTime).
		SetTotal(receipt.Total).
		Save(ctx)
}

func (r *ReceiptRepositoryImpl) DeleteReceipt(ctx context.Context, id int) error {
	return r.client.Receipt.DeleteOneID(id).Exec(ctx)
}
