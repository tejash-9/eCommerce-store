package internal

import (
	"os"
	"log"
	"fmt"
	"sync"
	"strconv"
)

type ShoppingEngine interface {
	RegisterUser(name string, email string) (*user, error)
	GetUserByUsername(username string) (*user, error)
}

type shoppingEngine struct {
	Users   			map[string]*user
	UserMap         	map[string]string
	Coupons         	map[string]string
	OrderThreshold  	int
	Inventory           *inventory
	OrderBook  			*orderBook
	CouponMutex			*sync.Mutex
}

func GetAppInstance() ShoppingEngine {
	instance.Do(func() {
		threshold, err := strconv.Atoi(os.Getenv(OrderThresholdEnv))
		if err != nil {
			log.Printf("Warning: Unable to get order threshold for discount, using default value. Error: %v", err)
			// Default order threshold for discount
            threshold = 30
		}
		shoppingApp = &shoppingEngine{
			Users: 			make(map[string]*user),
			UserMap: 		make(map[string]string),
			OrderBook:      &orderBook{
				Counter: 1,
			},
			Inventory:      &inventory{},
			OrderThreshold: threshold,
			CouponMutex: 	&sync.Mutex{},
		}
	})
	return shoppingApp
}

func (s *shoppingEngine) AddToCart(userId string, productId string, quantity int) (map[string]int, error) {
	// Check if user exists
	_, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}

	// Initialize the user's cart if it doesn't exist yet
    if s.Users[userId].Cart == nil {
        s.Users[userId].Cart = make(map[string]int)
    }
    
    // Add or update the product quantity
    s.Users[userId].Cart[productId] += quantity

	return s.Users[userId].Cart, nil
}

func (s *shoppingEngine) GetDiscountCoupon(userId string) (string, error) {
	// Check if user exists
	_, err := s.GetUser(userId)
	if err != nil {
		return "", err
	}

	if s.OrderBook.Counter % s.OrderThreshold != 0 {
		return "", fmt.Errorf("Dicount code not applicable!")
	}
	if s.Coupons[userId] == "" {
		s.Coupons[userId] = generateCouponCode(5)
	}
	return s.Coupons[userId], nil
}

func (s *shoppingEngine) Checkout(userId string, couponCode string) (*order, error) {
	_, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}

	if s.Users[userId].Cart == nil {
		return nil, fmt.Errorf("Cart is empty!")
	}

	var amount, discount float64
	for key, value := range s.Users[userId].Cart {
		amount += float64(float64(s.Products[key].Price) * float64(value))
	}
	var currentOrder *order

	if couponCode != "" && s.OrderBook.Counter % s.OrderThreshold == 0 {
		if s.Coupons[userId] == "" || s.Coupons[userId] != couponCode {
			return nil, fmt.Errorf("Invalid coupon code!")
		}

		s.CouponMutex.Lock()
		defer s.CouponMutex.Unlock()

		if s.OrderBook.Counter % s.OrderThreshold != 0 {
			return nil, fmt.Errorf("Coupon has expired!")
		}
        // Apply the 10% discount to the entire cart
        discount := amount * 0.10

		// Process and validate the cart
		var err error
		currentOrder, err = s.PlaceOrder(userId, amount, couponCode, discount)
		if err != nil {
			return nil, fmt.Errorf("Error processing order: %v", err)
		}
        delete(s.Coupons, userId)
	} else {
		// Process and validate the cart
		var err error
		currentOrder, err = s.PlaceOrder(userId, amount, couponCode, discount)
		if err != nil {
			return nil, fmt.Errorf("Error processing order: %v", err)
		}
	}
	// Clear the cart after checkout
	s.Users[userId].Cart = nil
	
	return currentOrder, nil
}