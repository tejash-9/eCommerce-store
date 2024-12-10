package internal

import (
	"fmt"
	"sync"
)

type OrderBook interface {
	TotalSoldItems()        int
	TotalPurchaseAmount()   float64
	ListDiscountCoupons()   []string
	TotalDiscountAmount()   float64
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

func newOrderBook() *orderBook {
	return &orderBook{
		Counter: 1,
		OrderMutex: &sync.Mutex{},
		Orders: make(map[string]*order), 
		OrdersByUserId: make(map[string][]*order),
		ItemsSold: 0,
		PurchaseAmount: 0.0,
		TotalDiscount: 0.0,
		AppliedCoupons: []string{},
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
			Logger.Sugar().Errorf("Coupon has expired for user: %s", userId)
			return nil, fmt.Errorf("Coupon has expired!") // Return error if coupon is expired
		}
		// Apply a 10% discount on the order total
		discount = amount * 0.10
	}

	var processedItems []string // List to track products whose stock was adjusted
	items := 0 					// Total number of items in the user's cart

	// Iterate through each item in the user's cart to adjust stock
	for key, value := range s.Users[userId].Cart {
		if !s.Inventory.Products[key].RemoveFromStock(value) {
			// Rollback any stock changes if a product is out of stock
			s.RollbackStock(userId, processedItems)
			Logger.Sugar().Errorf("Product %s is out of stock", key)
			return nil, fmt.Errorf("product %s is out of stock", key) // Return error if stock is insufficient
		}
		// Track processed items for rollback if needed
		processedItems = append(processedItems, key)
		items += value
	}

	// Calculate the final amount after applying the discount
	finalAmount := amount - discount

	// Update the order book with the new order data
	s.OrderBook.ItemsSold += items            // Increment total items sold
	s.OrderBook.PurchaseAmount += finalAmount // Increment total purchase amount
	s.OrderBook.TotalDiscount += discount     // Increment total discount applied

	// Generate a new unique order ID and create the order object
	id := generateUUID()
	order := newOrder(id, userId, s.Users[userId].Cart, amount, coupon, discount, finalAmount)

	// Store the newly created order in the order book
	s.OrderBook.Orders[id] = order
	// Associate the order with the user
	s.OrderBook.OrdersByUserId[userId] = append(s.OrderBook.OrdersByUserId[userId], order)
	s.OrderBook.Counter++ // Increment the order counter

	// Clear the user's cart after placing the order
	s.Users[userId].Cart = nil

	Logger.Sugar().Infof("Order placed successfully! total amount: %f, total discount: %f, item count: %d", s.OrderBook.PurchaseAmount, s.OrderBook.TotalDiscount, s.OrderBook.ItemsSold)
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

// TotalSoldItems returns the total number of items sold in the order book
func (o *orderBook) TotalSoldItems() int {
	return o.ItemsSold
}

// TotalPurchaseAmount returns the total amount of money spent on purchases in the order book
func (o *orderBook) TotalPurchaseAmount() float64 {
	return o.PurchaseAmount
}

// ListDiscountCoupons returns a list of all applied discount coupons
func (o *orderBook) ListDiscountCoupons() []string {
	return o.AppliedCoupons
}

// TotalDiscountAmount returns the total amount of discount applied across all orders
func (o *orderBook) TotalDiscountAmount() float64 {
	return o.TotalDiscount
}