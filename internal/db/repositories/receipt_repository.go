package repositories

import (
	"billing/internal/db/ent"
	"billing/internal/db/ent/receipt"
	"context"
	"log"
)

// ReceiptRepository defines the methods for interacting with receipts, including transaction support.
type ReceiptRepository interface {
	CreateReceipt(ctx context.Context, receipt *ent.Receipt) (*ent.Receipt, error)
	CreateReceiptWithItems(ctx context.Context, receipt *ent.Receipt, items []*ent.Item) (*ent.Receipt, error)
	GetReceiptByID(ctx context.Context, id int) (*ent.Receipt, error)
	UpdateReceipt(ctx context.Context, receipt *ent.Receipt) (*ent.Receipt, error)
	DeleteReceipt(ctx context.Context, id int) error
	BeginTx(ctx context.Context) (*ent.Tx, error)
	CommitTx(tx *ent.Tx) error
	RollbackTx(tx *ent.Tx) error
}

// receiptRepository is the concrete implementation of ReceiptRepository.
type receiptRepository struct {
	client *ent.Client
}

// NewReceiptRepository creates a new instance of ReceiptRepository.
func NewReceiptRepository(client *ent.Client) ReceiptRepository {
	return &receiptRepository{client: client}
}

// rollbackOnError handles transaction rollback and logs any errors that occur during rollback
func (r *receiptRepository) rollbackOnError(tx *ent.Tx, err error) {
	if rollbackErr := r.RollbackTx(tx); rollbackErr != nil {
		log.Printf("failed to rollback transaction: %v, original error: %v", rollbackErr, err)
	}
}

func (r *receiptRepository) CreateReceipt(ctx context.Context, receipt *ent.Receipt) (*ent.Receipt,
	error) {
	return r.client.Receipt.Create().
		SetRetailer(receipt.Retailer).
		SetPurchaseDate(receipt.PurchaseDate).
		SetPurchaseTime(receipt.PurchaseTime).
		SetTotal(receipt.Total).
		SetPoints(receipt.Points).
		Save(ctx)
}

func (r *receiptRepository) GetReceiptByID(ctx context.Context, id int) (*ent.Receipt, error) {
	return r.client.Receipt.Query().
		Where(receipt.ID(id)).
		WithItems().
		Only(ctx)
}

func (r *receiptRepository) UpdateReceipt(ctx context.Context, receipt *ent.Receipt) (*ent.Receipt, error) {
	return r.client.Receipt.UpdateOneID(receipt.ID).
		SetRetailer(receipt.Retailer).
		SetPurchaseDate(receipt.PurchaseDate).
		SetPurchaseTime(receipt.PurchaseTime).
		SetTotal(receipt.Total).
		SetPoints(receipt.Points).
		Save(ctx)
}

func (r *receiptRepository) DeleteReceipt(ctx context.Context, id int) error {
	return r.client.Receipt.DeleteOneID(id).Exec(ctx)
}

func (r *receiptRepository) BeginTx(ctx context.Context) (*ent.Tx, error) {
	return r.client.Tx(ctx)
}

func (r *receiptRepository) CommitTx(tx *ent.Tx) error {
	return tx.Commit()
}

func (r *receiptRepository) RollbackTx(tx *ent.Tx) error {
	return tx.Rollback()
}

// CreateReceiptWithItems creates a new receipt with its associated items in a single transaction.
// If any part of the operation fails, the entire transaction is rolled back.
func (r *receiptRepository) CreateReceiptWithItems(ctx context.Context, receipt *ent.Receipt, items []*ent.Item) (*ent.Receipt, error) {
	tx, err := r.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	createdReceipt, err := tx.Receipt.Create().
		SetRetailer(receipt.Retailer).
		SetPurchaseDate(receipt.PurchaseDate).
		SetPurchaseTime(receipt.PurchaseTime).
		SetTotal(receipt.Total).
		SetPoints(receipt.Points).
		Save(ctx)
	if err != nil {
		r.rollbackOnError(tx, err)
		return nil, err
	}

	for _, item := range items {
		if _, err := tx.Item.Create().
			SetShortDescription(item.ShortDescription).
			SetPrice(item.Price).
			SetReceiptsID(createdReceipt.ID).
			Save(ctx); err != nil {
			r.rollbackOnError(tx, err)
			return nil, err
		}
	}

	if err := r.CommitTx(tx); err != nil {
		return nil, err
	}

	return createdReceipt, nil
}
