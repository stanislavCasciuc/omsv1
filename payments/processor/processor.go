package processor

import (
	pb "github.com/stanislavCasciuc/common/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
