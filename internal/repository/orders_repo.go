package repository

import (
	"context"
	"fmt"

	"github.com/bekzourdk/orders/internal/domain/order"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
)

type OrdersRepository struct {
	con *pgx.Conn
}

func NewOrdersRepository(conn *pgx.Conn) *OrdersRepository {
	return &OrdersRepository{
		con: conn,
	}
}

func (u *OrdersRepository) Create(ctx context.Context, order *order.Order) error {
	if err := u.con.QueryRow(
		ctx,
		createOrderQuery,
		order.Details,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(
		&order.ID,
	); err != nil {
		return fmt.Errorf("OrdersRepository.Create: %w", err)
	}

	return nil
}

func (u *OrdersRepository) GetAll(ctx context.Context) ([]*order.Order, error) {
	rows, err := u.con.Query(ctx, selectAll)
	if err != nil {
		return nil, fmt.Errorf("OrdersRepository.GetAll: %w", err)
	}

	orders := []*order.Order{}

	for rows.Next() {
		var order order.Order

		if err := rows.Scan(
			&order.ID,
			&order.Details,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("OrdersRepository.GetAll: %w", err)
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func (u *OrdersRepository) FindByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	var order order.Order

	if err := u.con.QueryRow(
		ctx,
		findByID,
		id,
	).Scan(
		&order.ID,
		&order.Details,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("OrdersRepository.FindByID: %w", err)
	}

	return &order, nil
}

func (u *OrdersRepository) UpdateByID(ctx context.Context, order *order.Order) error {
	res, err := u.con.Exec(
		ctx,
		updateByID,
		order.Details,
		order.UpdatedAt,
		order.ID,
	)
	if err != nil {
		return fmt.Errorf("OrdersRepository.UpdateByID: %w", err)
	}

	if res.RowsAffected() < 1 {
		return fmt.Errorf("OrdersRepository.UpdateByID: no rows affected")
	}

	return nil
}

func (u *OrdersRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	res, err := u.con.Exec(ctx, deleteByID, id)
	if err != nil {
		return fmt.Errorf("OrdersRepository.DeleteByID: %w", err)
	}

	if res.RowsAffected() < 1 {
		return fmt.Errorf("OrdersRepository.DeleteByID: no rows affected")
	}

	return nil
}
