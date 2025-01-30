package repositories

import (
	"billing/internal/db/ent"
	"context"
)

// ItemRepository defines the methods for interacting with items.
type ItemRepository interface {
	CreateItem(ctx context.Context, item *ent.Item) (*ent.Item, error)
	GetItemByID(ctx context.Context, id int) (*ent.Item, error)
	UpdateItem(ctx context.Context, item *ent.Item) (*ent.Item, error)
	DeleteItem(ctx context.Context, id int) error
}

// ItemRepositoryImpl is the concrete implementation of ItemRepository.
type ItemRepositoryImpl struct {
	client *ent.Client
}

// NewItemRepository creates a new instance of ItemRepositoryImpl.
func NewItemRepository(client *ent.Client) ItemRepository {
	return &ItemRepositoryImpl{client: client}
}

func (r *ItemRepositoryImpl) CreateItem(ctx context.Context, item *ent.Item) (*ent.Item, error) {
	return r.client.Item.Create().
		SetShortDescription(item.ShortDescription).
		SetPrice(item.Price).
		Save(ctx)
}

func (r *ItemRepositoryImpl) GetItemByID(ctx context.Context, id int) (*ent.Item, error) {
	return r.client.Item.Get(ctx, id)
}

func (r *ItemRepositoryImpl) UpdateItem(ctx context.Context, item *ent.Item) (*ent.Item, error) {
	return r.client.Item.UpdateOneID(item.ID).
		SetShortDescription(item.ShortDescription).
		SetPrice(item.Price).
		Save(ctx)
}

func (r *ItemRepositoryImpl) DeleteItem(ctx context.Context, id int) error {
	return r.client.Item.DeleteOneID(id).Exec(ctx)
}
