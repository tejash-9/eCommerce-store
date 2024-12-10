package internal

import (
	"fmt"
	"sync"
)

type order struct {
	Id              string           	`json:"id"`
	UserId          string          	`json:"user_id"`
	OrderCart       map[string]int  	`json:"order_cart"`
	CartTotal     	float64             `json:"amount"`
	Discount        float64            	`json:"discount"`
	DiscountCoupon  string            	`json:"discount_coupon"`
	AmountToPay     float64          	`json:"amount_to_pay"`
}

func newOrder(id string, userId string, cart map[string]int, amount float64, coupon string, discount float64, finalAmount float64) *order {
	return &order{
		Id:             id,
		UserId:         userId,
		OrderCart:      cart,
		CartTotal:      amount,
		DiscountCoupon: coupon,
		Discount:       discount,
		AmountToPay:    finalAmount,
	}
}

func (s *shoppingEngine) PlaceOrder(userId string, amount float64, coupon string, discount float64) (*order, error) {
	items, err := s.ValidateCart(userId)
	if err != nil {
		return nil, err
	}

	// Total payable amount 
	finalAmount := amount - discount

	// Update the total items, purchase and discount amount
	s.OrderBook.ItemsSold += *items
	s.OrderBook.PurchaseAmount += finalAmount
	s.OrderBook.TotalDiscount += discount

	if coupon != "" {
		s.OrderBook.AppliedCoupons = append(s.OrderBook.AppliedCoupons, coupon)
	}
	// Generate a new order ID and create the order
	id := generateUUID()
	order := newOrder(id, userId, s.Users[userId].Cart, amount, coupon, discount, finalAmount)

	// Store the order in the orderBook
	s.OrderBook.Orders[id] = order
	s.OrderBook.OrdersByUserId[userId] = append(s.OrderBook.OrdersByUserId[userId], order)
	s.OrderBook.Counter++

	return order, nil
}