package stripe

import (
	"fmt"
	"log"

	"github.com/stanislavCasciuc/common"
	pb "github.com/stanislavCasciuc/common/api"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

var gatewayHTTPAddr = common.EnvString("GATEWAY_HTTP_ADDR", "http://localhost:8080")

type Stripe struct{}

func NewProcessor() *Stripe {
	return &Stripe{}
}

func (s *Stripe) CreatePaymentLink(o *pb.Order) (string, error) {
	log.Printf("Creating payment link for order %v", o)

	gatewaySuccessURL := fmt.Sprintf("%s/success.html", gatewayHTTPAddr)

	items := []*stripe.CheckoutSessionLineItemParams{}
	for _, item := range o.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PriceID), //stripe.String(item.PriceID)
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessURL),
	}
	result, err := session.New(params)
	if err != nil {
		return "", nil
	}

	return result.URL, nil
}
