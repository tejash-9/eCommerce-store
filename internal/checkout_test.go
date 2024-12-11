package internal

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a mock shopping engine and users
func createMockEngine() *shoppingEngine {
	engine := &shoppingEngine{
		Users:     make(map[string]*user),
		UserMap:   make(map[string]string),
		Coupons:   make(map[string]string),
		Inventory: newInventory(),
		OrderBook: newOrderBook(),
		DiscountInterval: 2, // Every 2nd order is applicable for discount
	}

	return engine
}

// Test AddToCart for a valid product
func TestAddToCart_Success(t *testing.T) {
	// Arrange
	shoppingApp := createMockEngine()

	// Register a seller and get the user ID
	seller, err := shoppingApp.RegisterUser("David", "david@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, seller)
	
	// Act
	p1, err := shoppingApp.RegisterProduct("Product 1", "Description of product 1", 10, seller.Id, 99.99)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	// Register a user and get the user ID
	user, err := shoppingApp.RegisterUser("Adity", "aditya@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)

	cart, err := shoppingApp.AddToCart(user.Id, p1.Id, 5)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 5, cart[p1.Id]) // Should be 2 (existing) + 3 (added)
}

// Test AddToCart for a non-existent user
func TestAddToCart_UserNotFound(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a seller and get the user ID
	seller, err := shoppingApp.RegisterUser("Shalom", "shalom@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, seller)
	
	// Act
	p1, err := shoppingApp.RegisterProduct("Product 1", "Description of product 1", 10, seller.Id, 99.99)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	cart, err := shoppingApp.AddToCart("NonExistentUser", p1.Id, 5)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, cart)
}

// Test AddToCart for a non-existent product
func TestAddToCart_ProductNotFound(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a user and get the user ID
	user, err := shoppingApp.RegisterUser("Aditya", "aditya@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)

	cart, err := shoppingApp.AddToCart(user.Id, "NonExistentProduct", 5)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, cart)
}

// Test GetDiscountCoupon for a user with a valid coupon
func TestGetDiscountCoupon_ValidCoupon(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a seller and get the user ID
	seller, err := shoppingApp.RegisterUser("Ram", "ram@example.com")

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

	// Act
	coupon, err := shoppingApp.GetDiscountCoupon(user.Id)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, coupon)
}

// Test GetDiscountCoupon for a user without a valid coupon
func TestGetDiscountCoupon_InvalidCoupon(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a seller and get the user ID
	seller, err := shoppingApp.RegisterUser("kiran", "kiran@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, seller)
	
	// Act
	p1, err := shoppingApp.RegisterProduct("Product 1", "Description of product 1", 10, seller.Id, 99.99)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	// Register a user and get the user ID
	user, err := shoppingApp.RegisterUser("Aditya", "aditya@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Act
	coupon, err := shoppingApp.GetDiscountCoupon(user.Id)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, coupon) // Coupon should not be generated
}

// Test PlaceOrder with a successful order
func TestCheckout_Success_NoCoupon(t *testing.T) {
	// Arrange
	shoppingApp := createMockEngine()

	// Register a seller and get the user ID
	seller, err := shoppingApp.RegisterUser("salman", "salman@example.com")

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

	_, err = shoppingApp.AddToCart(user.Id, p2.Id, 5)

	// Assert
	assert.NoError(t, err)

	// Act
	order, err := shoppingApp.Checkout(user.Id, "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
}

// Test Checkout for a valid order with a coupon
func TestCheckout_Success_WithCoupon(t *testing.T) {
	shoppingApp := createMockEngine()
	shoppingApp.DiscountInterval = 2

	// Register a seller and get the user ID
	seller, err := shoppingApp.RegisterUser("Devi prasad", "devi@example.com")

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
	user, err := shoppingApp.RegisterUser("Adity", "aditya@example.com")

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

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
}

// Test Checkout for a user with an invalid coupon
func TestCheckout_InvalidCoupon(t *testing.T) {
	shoppingApp := createMockEngine()
	shoppingApp.DiscountInterval = 2

	// Register a seller and get the user ID
	seller, err := shoppingApp.RegisterUser("Piyush", "piyush@example.com")

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
	order, err := shoppingApp.Checkout(user.Id, "INVALID_COUPON")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, order)
}

// Test Checkout for an empty cart
func TestCheckout_EmptyCart(t *testing.T) {
	shoppingApp := createMockEngine()
	shoppingApp.DiscountInterval = 2

	// Register a seller and get the user ID
	seller, err := shoppingApp.RegisterUser("John", "john@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, seller)
	
	// Act
	p1, err := shoppingApp.RegisterProduct("Product 1", "Description of product 1", 10, seller.Id, 99.99)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	// Register a user and get the user ID
	user, err := shoppingApp.RegisterUser("Aditya", "aditya@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)

	// Act
	order, err := shoppingApp.Checkout(user.Id, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, order)
}

// Test Checkout with insufficient stock
func TestCheckout_InsufficientStock(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a seller and get the user ID
	seller, err := shoppingApp.RegisterUser("John Doe", "john.doe@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, seller)
	
	// Act
	p1, err := shoppingApp.RegisterProduct("Product 1", "Description of product 1", 10, seller.Id, 99.99)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	// Register a user and get the user ID
	user, err := shoppingApp.RegisterUser("Aditya", "aditya@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)

	cart, err := shoppingApp.AddToCart(user.Id, p1.Id, 12)

	// Assert
	assert.NoError(t, err)

	// Act
	order, err := shoppingApp.Checkout(user.Id, "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, user.Cart, cart) // Cart should remain unchanged
}