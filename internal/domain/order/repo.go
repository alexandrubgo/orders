package order

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type Repo interface {
	Create(ctx context.Context, order *Order) error
	GetAll(ctx context.Context) ([]*Order, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)
	UpdateByID(ctx context.Context, user *Order) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
