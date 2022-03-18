package order

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Order struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Details   string    `json:"details,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewOrder(details string) *Order {
	return &Order{
		Details:   details,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (o *Order) Update() {
	o.UpdatedAt = time.Now()
}
