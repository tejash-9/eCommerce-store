package internal

import (
	"fmt"
)

// AddToCart adds a product to the user's cart
func (s *shoppingEngine) AddToCart(userId string, productId string, quantity int) (map[string]int, error) {
	// Check if user exists
	_, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}

	// check if product exists
	_, err = s.GetProduct(productId)
	if err != nil {
		return nil, err
	}

	// Add or update the product quantity
	s.Users[userId].Cart[productId] += quantity

	Logger.Sugar().Infof("Product %s added to cart successfully by user: %s", productId, userId)
	return s.Users[userId].Cart, nil
}

func (s *shoppingEngine) GetCart(userId string) (map[string]int, error) {
	// Check if user exists
	_, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}
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
		return "", fmt.Errorf("Discount code not applicable")
	}

	return s.GenerateDiscountCouponForUser(userId), nil
}

// Checkout processes the user's cart and applies a coupon if valid
func (s *shoppingEngine) Checkout(userId string, couponCode string) (*order, error) {
	// Check if user exists
	_, err := s.GetUser(userId)
	if err != nil {
		return nil, err
	}

	// Ensure cart is not empty
	if len(s.Users[userId].Cart) == 0 {
		return nil, fmt.Errorf("Cart is empty")
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
	} else {
		// Place the order without any discount
		currentOrder, err = s.PlaceOrder(userId, amount, couponCode)
		if err != nil {
			return nil, err
		}
	}
	s.Coupons = make(map[string]string)
	Logger.Sugar().Info("Checkout successful!")
	return currentOrder, nil
}

func (s *shoppingEngine) GenerateDiscountCouponForUser(userId string) string {
	// Generate a new coupon code if it doesn't already exist
	if s.Coupons[userId] == "" {
		s.Coupons[userId] = generateCouponCode(5)
		Logger.Sugar().Info("Discount coupon generated successfully!")
	}
	return s.Coupons[userId]
}