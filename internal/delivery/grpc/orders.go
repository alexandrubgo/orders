package grpc

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bekzourdk/orders/internal/domain/order"
	"github.com/bekzourdk/orders/internal/pb"
	"github.com/bekzourdk/orders/internal/service"
)

type OrdersGRPC struct {
	ordersService service.OrdersService
	pb.UnimplementedOrdersServiceServer
}

func NewOrdersGRPC(
	ordersService service.OrdersService,
) *OrdersGRPC {
	return &OrdersGRPC{
		ordersService: ordersService,
	}
}

func (u *OrdersGRPC) GetAll(ctx context.Context, req *pb.GetAllOrdersRequests) (*pb.GetAllOrdersResponse, error) {
	res, err := u.ordersService.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("OrdersGRPC.GetAll: %w", err)
	}

	ordersProto := make([]*pb.Order, len(res))

	for i, v := range res {
		ordersProto[i] = OrderDomainToProto(v)
	}

	return &pb.GetAllOrdersResponse{
		Orders: ordersProto,
	}, nil
}

func (p *OrdersGRPC) Create(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	orderCandidate := order.NewOrder(req.Details)

	err := p.ordersService.Create(ctx, orderCandidate)
	if err != nil {
		return nil, fmt.Errorf("OrdersGRPC.Create: %w", err)
	}

	return &pb.CreateOrderResponse{
		Order: OrderDomainToProto(orderCandidate),
	}, nil
}

func (p *OrdersGRPC) FindByID(ctx context.Context, req *pb.FindByIDRequests) (*pb.FindByIDResponse, error) {
	id := uuid.FromStringOrNil(req.Id)
	if id == uuid.Nil {
		return nil, fmt.Errorf("id is invalid")
	}

	res, err := p.ordersService.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("OrdersGRPC.FindByID: %w", err)
	}

	return &pb.FindByIDResponse{
		Order: OrderDomainToProto(res),
	}, nil
}
func (p *OrdersGRPC) UpdateByID(ctx context.Context, req *pb.UpdateByIDRequests) (*pb.UpdateByIDResponse, error) {
	order := OrderProtoToDomain(req.Order)

	err := p.ordersService.UpdateByID(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("OrdersGRPC.UpdateByID: %w", err)
	}

	return &pb.UpdateByIDResponse{
		Order: OrderDomainToProto(order),
	}, nil
}
func (p *OrdersGRPC) DeleteByID(ctx context.Context, req *pb.DeleteByIDRequests) (*pb.DeleteByIDResponse, error) {
	id := uuid.FromStringOrNil(req.Id)
	if id == uuid.Nil {
		return nil, fmt.Errorf("id is invalid")
	}

	err := p.ordersService.DeleteByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("OrdersGRPC.DeleteByID: %w", err)
	}

	return &pb.DeleteByIDResponse{
		Id: id.String(),
	}, nil
}

func OrderDomainToProto(order *order.Order) *pb.Order {
	return &pb.Order{
		Id:        order.ID.String(),
		Details:   order.Details,
		CreatedAt: timestamppb.New(order.CreatedAt),
		UpdatedAt: timestamppb.New(order.UpdatedAt),
	}
}

func OrderProtoToDomain(orderProto *pb.Order) *order.Order {
	return &order.Order{
		ID:        uuid.FromStringOrNil(orderProto.Id),
		Details:   orderProto.Details,
		CreatedAt: orderProto.CreatedAt.AsTime(),
		UpdatedAt: orderProto.UpdatedAt.AsTime(),
	}
}
