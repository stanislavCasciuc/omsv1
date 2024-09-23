package main

import (
	"context"
	"testing"

	"github.com/stanislavCasciuc/common/api"
	"github.com/stanislavCasciuc/payments/processor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()

	svc := NewService(processor)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error = %v, wants nil", err)

		}

		if link == "" {
			t.Error("CreatePayment() link is empty")
		}
	})
}
