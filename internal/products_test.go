package internal

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// Test RegisterProduct with a valid seller and product details
func TestRegisterProduct_Success(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a seller user
	seller, _ := shoppingApp.RegisterUser("Seller", "seller@example.com")

	// Act
	product, err := shoppingApp.RegisterProduct("Product 1", "Description of product 1", 10, seller.Id, 99.99)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Product 1", product.Name)
	assert.Equal(t, 10, product.Quantity)
	assert.Equal(t, 99.99, product.Price)
	assert.Len(t, shoppingApp.Inventory.Products, 1)
	assert.Len(t, shoppingApp.Inventory.ProductsBySeller[seller.Id], 1)
}

// Test RegisterProduct with an invalid seller (user doesn't exist)
func TestRegisterProduct_SellerNotFound(t *testing.T) {
	shoppingApp := createMockEngine()

	// Act
	product, err := shoppingApp.RegisterProduct("Product 1", "Description of product 1", 10, "nonExistentSeller", 99.99)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, product)
}

// Test RegisterProduct with an existing product name by the same seller
func TestRegisterProduct_ProductAlreadyExists(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a seller user
	seller, _ := shoppingApp.RegisterUser("Seller", "seller@example.com")

	// Register the first product
	p1, err := shoppingApp.RegisterProduct("Product 1", "Description of product 1", 10, seller.Id, 99.99)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, p1)

	// Act
	p2, err := shoppingApp.RegisterProduct("Product 1", "Another description", 20, seller.Id, 89.99)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, p2)
	assert.Equal(t, "product with name already exists by the seller", err.Error())
}

// Test GetProduct with an existing product ID
func TestGetProduct_Success(t *testing.T) {
	shoppingApp := createMockEngine()

	// Register a seller user
	seller, _ := shoppingApp.RegisterUser("Seller", "seller@example.com")

	// Register a product
	product, _ := shoppingApp.RegisterProduct("Product 1", "Description of product 1", 10, seller.Id, 99.99)

	// Act
	retrievedProduct, err := shoppingApp.GetProduct(product.Id)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, product, retrievedProduct)
}

// Test GetProduct with a non-existent product ID
func TestGetProduct_ProductNotFound(t *testing.T) {
	shoppingApp := createMockEngine()

	// Act
	retrievedProduct, err := shoppingApp.GetProduct("nonExistentProductId")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, retrievedProduct)
	assert.Equal(t, "Product not found", err.Error())
}
