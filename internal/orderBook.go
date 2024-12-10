package internal

import (
	"fmt"
	"sync"
)

type OrderBook interface {
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

func (s *shoppingEngine) OrderHistory() OrderBook {
	return s.OrderBook
}

func (s *shoppingEngine) ValidateCart(userId string) (*int, error) {
	s.OrderBook.OrderMutex.Lock()
	defer s.OrderBook.OrderMutex.Unlock()

	items := 0
	var processedItems []string
	for key, value := range s.Users[userId].Cart {
		// Attempt to remove the product from stock
		if !s.Inventory.Products[key].RemoveFromStock(value) {
			// Rollback previously removed products and return error
			s.RollbackStock(userId, processedItems)
			return nil, fmt.Errorf("product %s is out of stock", key)
		}
		processedItems = append(processedItems, key)
		items += value
	}
	return &items, nil
}

func (s *shoppingEngine) RollbackStock(userId string, products []string) {
	for _, productId := range products {
		if s.Inventory.Products[productId] != nil {
			s.Inventory.Products[productId].AddToStock(s.Users[userId].Cart[productId])
		}
	}
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