package internal

import (
	"os"
	"fmt"
	"strconv"
)

type ShoppingEngine interface {
	RegisterUser(name string, email string) (*user, error)
	GetUser(userId string) (*user, error)
	GetUserByUsername(username string) (*user, error)
	RegisterProduct(name string, description string, quantity int, sellerId string, price float64) (*product, error)
	GetProduct(productId string) (*product, error)
	AddToCart(userId string, productId string, quantity int) (map[string]int, error)
	GetDiscountCoupon(userId string) (string, error)
	Checkout(userId string, couponCode string) (*order, error)
	OrderHistory() OrderBook
}

type shoppingEngine struct {
	Users             map[string]*user         // Map of users, indexed by userId
	UserMap           map[string]string        // Map to store username mappings
	Coupons           map[string]string        // Coupons by userId
	DiscountInterval  int                      // Discount interval (every N orders)
	Inventory         *inventory               // Inventory system with products
	OrderBook         *orderBook               // Order history tracking
}

// GetAppInstance creates and returns a singleton instance of the ShoppingEngine
func GetAppInstance() ShoppingEngine {
	instance.Do(func() {
		// Default discount interval
		interval, err := strconv.Atoi(os.Getenv("DISCOUNT_INTERVAL"))
		if err != nil {
			Logger.Sugar().Debugf("Unable to get discount interval, using default value: %v", err)
			interval = 10 // Default discount interval if not provided
		}
		orderBook := newOrderBook()
		inventory := newInventory()
		// Initialize the shopping engine with default values
		shoppingApp = &shoppingEngine{
			Users:            make(map[string]*user),
			UserMap:          make(map[string]string),
			OrderBook:        orderBook,
			Inventory:        inventory,
			DiscountInterval: interval,
			Coupons:          make(map[string]string),
		}
	})
	return shoppingApp
}

// AddToCart adds a product to the user's cart
func (s *shoppingEngine) AddToCart(userId string, productId string, quantity int) (map[string]int, error) {
	// Check if user exists
	_, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}

	// Initialize the user's cart if it doesn't exist
	if s.Users[userId].Cart == nil {
		s.Users[userId].Cart = make(map[string]int)
	}

	// Add or update the product quantity
	s.Users[userId].Cart[productId] += quantity

	Logger.Sugar().Infof("Product %s added to cart successfully", productId)
	return s.Users[userId].Cart, nil
}

// GetDiscountCoupon retrieves a discount coupon for the user, if applicable
func (s *shoppingEngine) GetDiscountCoupon(userId string) (string, error) {
	// Check if user exists
	_, err := s.GetUser(userId)
	if err != nil {
		return "", err
	}

	// Check if the order counter matches the discount interval
	if s.OrderBook.Counter%s.DiscountInterval != 0 {
		return "", fmt.Errorf("Discount code not applicable!")
	}

	// Generate a new coupon code if it doesn't already exist
	if s.Coupons[userId] == "" {
		s.Coupons[userId] = generateCouponCode(5)
	}

	Logger.Sugar().Info("Discount coupon generated successfully!")
	return s.Coupons[userId], nil
}

// Checkout processes the user's cart and applies a coupon if valid
func (s *shoppingEngine) Checkout(userId string, couponCode string) (*order, error) {
	// Check if user exists
	_, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}

	// Ensure cart is not empty
	if s.Users[userId].Cart == nil {
		return nil, fmt.Errorf("Cart is empty!")
	}

	// Calculate total amount of items in the cart
	var amount float64
	for productId, quantity := range s.Users[userId].Cart {
		amount += float64(s.Inventory.Products[productId].GetPrice()) * float64(quantity)
	}

	var currentOrder *order
	// Check if coupon code is provided
	if couponCode != "" {
		// Validate the coupon code
		if s.Coupons[userId] == "" || s.Coupons[userId] != couponCode {
			return nil, fmt.Errorf("Invalid coupon code!")
		}

		// Place the order with discount
		currentOrder, err = s.PlaceOrder(userId, amount, couponCode)
		if err != nil {
			return nil, err
		}
		// Clear the coupon for the user after usage
		delete(s.Coupons, userId)
	} else {
		// Place the order without any discount
		currentOrder, err = s.PlaceOrder(userId, amount, couponCode)
		if err != nil {
			return nil, err
		}
	}
	Logger.Sugar().Info("Checkout successful!")
	return currentOrder, nil
}