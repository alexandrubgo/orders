package service

import (
	"context"
	"fmt"

	"github.com/bekzourdk/orders/internal/domain/order"
	uuid "github.com/satori/go.uuid"
)

type OrdersService struct {
	ordersRepository order.Repo
}

func NewOrdersService(
	ordersRepository order.Repo,
) *OrdersService {
	return &OrdersService{
		ordersRepository: ordersRepository,
	}
}

func (s *OrdersService) Create(ctx context.Context, order *order.Order) error {
	if err := s.ordersRepository.Create(ctx, order); err != nil {
		return fmt.Errorf("OrdersService.Create: %w", err)
	}

	return nil
}

func (s *OrdersService) GetAll(ctx context.Context) ([]*order.Order, error) {
	res, err := s.ordersRepository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("OrdersService.GetAll: %w", err)

	}

	return res, nil
}

func (s *OrdersService) FindByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	res, err := s.ordersRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("OrdersService.FindByID: %w", err)

	}

	return res, nil
}

func (s *OrdersService) UpdateByID(ctx context.Context, order *order.Order) error {
	order.Update()

	if err := s.ordersRepository.UpdateByID(ctx, order); err != nil {
		return fmt.Errorf("OrdersService.UpdateByID: %w", err)
	}

	return nil
}

func (s *OrdersService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	if err := s.ordersRepository.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("OrdersService.DeleteByID: %w", err)
	}

	return nil
}
