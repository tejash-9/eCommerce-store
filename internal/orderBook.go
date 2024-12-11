package internal

import (
	"fmt"
	"sync"
)

type OrderBook interface {
	GetAnalytics() (int, float64, float64, []string)
}

type orderBook struct {
	ItemsSold         int               	// Total number of items sold
	PurchaseAmount    float64           	// Total amount spent on all purchases
	TotalDiscount     float64           	// Total discount amount applied
	AppliedCoupons    []string          	// List of applied coupon codes
	Orders            map[string]*order 	// Map of all orders by orderId
	OrdersByUserId    map[string][]*order 	// Map of orders by userId
	OrderMutex        *sync.Mutex       	// Mutex to prevent race conditions in order history
	Counter           int               	// Counter for order numbering
}

// newOrderBook Creates and returns a new orderBook object
func newOrderBook() *orderBook {
	return &orderBook{
		Counter: 1,
		OrderMutex: &sync.Mutex{},
		Orders: make(map[string]*order), 
		OrdersByUserId: make(map[string][]*order),
	}
}

// OrderHistory returns the current order book interface
func (s *shoppingEngine) OrderHistory() OrderBook {
	return s.OrderBook
}

// PlaceOrder processes the order by adjusting inventory, updating the order book, and creating a new order
func (s *shoppingEngine) PlaceOrder(userId string, amount float64, coupon string) (*order, error) {
	// Lock to ensure thread-safe operations on inventory and order history
	s.OrderBook.OrderMutex.Lock()
	defer s.OrderBook.OrderMutex.Unlock()

	var discount float64
	// Apply coupon discount if a coupon is provided
	if coupon != "" {
		// Check if the coupon is valid based on a discount interval (e.g., every N orders)
		if s.OrderBook.Counter % s.DiscountInterval != 0 {
			Logger.Sugar().Debugf("Coupon has expired for user: %s", userId)
			return nil, fmt.Errorf("Coupon has expired")
		}

		if s.Coupons[userId] != coupon {
			return nil, fmt.Errorf("Invalid coupon code")
		}
		// Apply a 10% discount on the order total
		discount = amount * 0.10
	}

	var processedItems []string
	var items int

	// Iterate through each item in the user's cart to adjust stock
	for key, value := range s.Users[userId].Cart {
		if !s.Inventory.Products[key].RemoveFromStock(value) {
			// Rollback any stock changes if a product is out of stock
			Logger.Sugar().Debugf("Product %s is out of stock, rolling back the cart changes!", key)
			s.RollbackStock(userId, processedItems)

			return nil, fmt.Errorf("Product %s is out of stock", key)
		}
		// Track processed items for rollback if needed
		processedItems = append(processedItems, key)
		items += value
	}

	// Calculate the final amount after applying the discount
	finalAmount := amount - discount

	// Generate a new unique order ID and create the order object
	id := generateUUID()
	order := newOrder(id, userId, s.Users[userId].Cart, amount, coupon, discount, finalAmount)

	// Store the newly created order in the order book
	s.OrderBook.Orders[id] = order
	s.OrderBook.OrdersByUserId[userId] = append(s.OrderBook.OrdersByUserId[userId], order)
	s.OrderBook.Counter++ // Increment the order counter

	// Clear the user's cart after placing the order
	s.Users[userId].Cart = make(map[string]int)

	// Update the order book
	s.OrderBook.ItemsSold += items            
	s.OrderBook.PurchaseAmount += finalAmount
	s.OrderBook.TotalDiscount += discount
	if coupon != "" {
		s.OrderBook.AppliedCoupons = append(s.OrderBook.AppliedCoupons, coupon)
	}

	Logger.Sugar().Infof("Order placed successfully with id: %s by user: %s", id, userId)
	return order, nil
}

// RollbackStock reverts the stock changes for the specified products in the user's cart
func (s *shoppingEngine) RollbackStock(userId string, products []string) {
	// Rollback all changes made during the cart validation process
	for _, productId := range products {
		// Add back the quantity of each product to the stock
		if s.Inventory.Products[productId] != nil {
			s.Inventory.Products[productId].AddToStock(s.Users[userId].Cart[productId])
		}
	}
}

// GetAnalytics returns the information about analytics related to orders
func (o *orderBook) GetAnalytics() (int, float64, float64, []string) {
	o.OrderMutex.Lock()
	defer o.OrderMutex.Unlock()

	return o.ItemsSold, o.PurchaseAmount, o.TotalDiscount, o.AppliedCoupons
}