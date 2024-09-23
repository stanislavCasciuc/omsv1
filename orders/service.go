package main

import (
	"context"
	"log"

	"github.com/stanislavCasciuc/common"
	pb "github.com/stanislavCasciuc/common/api"
)

type service struct {
	store OrderStore
}

func NewService(store OrderStore) *service {
	return &service{store}
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	items, err := s.ValidateOrders(ctx, p)
	if err != nil {
		return &pb.Order{}, err
	}

	o := &pb.Order{
		ID:         "42",
		Status:     "pending",
		CustomerID: p.CustomerID,
		Items:      items,
	}
	return o, nil
}

func (s *service) ValidateOrders(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}

	mergedItems := mergeItemsQuantities(p.Items)
	log.Print(mergedItems)

	// validate with stock service

	//Temporary
	var itemsWithPrice []*pb.Item
	for _, i := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PriceID:  "price_1PvIGVLsu5wxhn3RDDfxCZSg",
			ID:       i.ID,
			Quantity: i.Quantity,
		})
	}

	return itemsWithPrice, nil
}

func mergeItemsQuantities(items []*pb.ItemWithQuantity) []*pb.ItemWithQuantity {
	merged := make([]*pb.ItemWithQuantity, 0)

	for _, item := range items {
		found := false
		for _, finalItem := range merged {
			if finalItem.ID == item.ID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
