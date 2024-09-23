package main

import (
	"context"

	pb "github.com/stanislavCasciuc/common/api"

)

type PaymentService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
