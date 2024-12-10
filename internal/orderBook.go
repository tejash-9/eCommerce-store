package internal

import (
	"fmt"
	"sync"
)

type OrderBook inter {
	TotalSoldItems() 		int
	TotalPurchaseAmount() 	float64
	ListDiscountCoupons() 	[]string 
	TotalDiscountAmount() 	float64
}

type orderBook struct {
	ItemsSold 			int
	PurchaseAmount     	float64
	TotalDiscount  		float64
	AppliedCoupons  	[]string
	Orders		        map[string]*order
	OrdersByUserId      map[string][]*order
	OrderMutex          *sync.Mutex
	Counter             int
}

func (s *shoppingEngine) OrderBook() OrderBook {
	return s.OrderBook
}

func (s *shoppingEngine) ValidateCart(userId string) (*int, error) {
	s.OrderBook.OrderMutex.Lock()
	defer s.OrderBook.OrderMutex.Unlock()

	items := 0
	for key, value := range s.Users[userId].Cart {
		if s.Products[key].Quantity < value {
			return nil, fmt.Errorf("product %s is out of stock")
		}
	}
	for key, value := range s.Users[userId].Cart {
		s.Products[key].Quantity = s.Products[key].Quantity - value
		items += value
	}
	return &items, nil
}

func (o *orderBook) TotalSoldItems() int {
	return o.ItemsSold
}

func (o *orderBook) TotalPurchaseAmount() float64 {
	return o.PurchaseAmount
}

func (o *orderBook) ListDiscountCoupons() []string {
	return o.AppliedCoupons
}

func (o *orderBook) TotalDiscountAmount() float64 {
	return o.TotalDiscount
}