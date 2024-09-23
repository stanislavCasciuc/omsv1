package main

import (
	"context"
	pb "github.com/stanislavCasciuc/common/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrders(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
}

type OrderStore interface {
	Create(context.Context) error
}
