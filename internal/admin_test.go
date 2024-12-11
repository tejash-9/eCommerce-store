package internal

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// Test TotalDiscountAmount
func TestAnalytics(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a seller and get the user ID
	seller, err := shoppingApp.RegisterUser("ken", "ken@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, seller)
	
	// Act
	p1, err := shoppingApp.RegisterProduct("Product 1", "Description of product 1", 10, seller.Id, 99.99)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	// Act
	p2, err := shoppingApp.RegisterProduct("Product 2", "Description of product 2", 20, seller.Id, 199.99)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, p2)

	// Register a user and get the user ID
	user, err := shoppingApp.RegisterUser("Aditya", "aditya@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)

	_, err = shoppingApp.AddToCart(user.Id, p1.Id, 5)

	// Assert
	assert.NoError(t, err)

	// Act
	order, err := shoppingApp.Checkout(user.Id, "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)

	amount1 := 99.99*5

	_, err = shoppingApp.AddToCart(user.Id, p2.Id, 5)

	// Assert
	assert.NoError(t, err)

	// Act
	coupon, err := shoppingApp.GetDiscountCoupon(user.Id)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, coupon)

	// Act
	order, err = shoppingApp.Checkout(user.Id, coupon)

	discount := 199.99*5*0.10
	amount2 := 199.99*5 - discount

	finalAmount := amount1 + amount2

	totalItems, totalAmount, totalDiscount, coupons := shoppingApp.OrderBook.GetAnalytics()
	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, 10, totalItems)
	assert.Equal(t, finalAmount, totalAmount)
	assert.Equal(t, discount, totalDiscount)
	assert.Equal(t, []string{coupon}, coupons)
}