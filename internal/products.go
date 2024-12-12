package internal

import (
	"fmt"
)

// product represents a single product in the inventory
type product struct {
	Id              string  `json:"id"` 			// Unique identifier for the product
	Name   			string  `json:"name"` 			// Name of the product
	Description    	string  `json:"description"` 	// Description of the product
	Quantity        int     `json:"quantity"` 		// Available stock quantity
	Price           float64 `json:"price"` 			// Price of the product
	SellerId        string  `json:"seller_id"` 		// Seller's unique identifier
}

// inventory manages the collection of products and their categorization by seller
type inventory struct {
	Products   			map[string]*product       // Mapping of product IDs to products
	ProductsBySeller 	map[string][]*product     // Mapping of seller IDs to their products
}

func newInventory() *inventory {
	return &inventory{
		Products: make(map[string]*product), 
		ProductsBySeller: make(map[string][]*product),
	}
}

// newProduct creates and returns a new product instance
func newProduct(id string, name string, description string, quantity int, seller_id string, price float64) *product {
	return &product{
		Id:             id,
		Name:           name,
		Description:    description,
		Quantity:       quantity,
		SellerId:       seller_id,
		Price:          price,
	}
}

// RegisterProduct adds a new product to the inventory if the seller is valid and the product doesn't already exist
func (s *shoppingEngine) RegisterProduct(name string, description string, quantity int, sellerId string, price float64) (*product, error) {
	// Ensure the seller is valid
	_, err := s.GetUser(sellerId)
	if err != nil {
		return nil, err // Return error if seller does not exist
	}

	// Check if product already exists for the seller
	for _, product := range s.Inventory.ProductsBySeller[sellerId] {
		if product.Name == name {
			return nil, fmt.Errorf("product with name already exists by the seller") // Error if the product already exists
		}
	}

	// Generate unique product ID and create the new product
	id := generateUUID()
	product := newProduct(id, name, description, quantity, sellerId, price)

	// Add product to the seller's inventory and global inventory
	s.Inventory.ProductsBySeller[sellerId] = append(s.Inventory.ProductsBySeller[sellerId], product)
	s.Inventory.Products[id] = product

	Logger.Sugar().Infof("Product %s registered successfully", id)
	return product, nil
}

// GetProduct fetches a product by its ID from the inventory
func (s *shoppingEngine) GetProduct(productId string) (*product, error) {
	if s.Inventory.Products[productId] == nil {
		return nil, fmt.Errorf("Product not found")
	}
	return s.Inventory.Products[productId], nil
}

func (s *shoppingEngine) RemoveProduct(productId string) error {
	if s.Inventory.Products[productId] == nil {
		return fmt.Errorf("Product not found")
	}
	s.Inventory.Products[productId] = nil
	return nil
}

// IsAvailable checks if the requested quantity of the product is in stock
func (p *product) IsAvailable(quantity int) bool {
	return p.Quantity >= quantity
}

// GetPrice returns the price of the product
func (p *product) GetPrice() float64 {
	return p.Price
}

// AddToStock increases the product's stock by the specified quantity
func (p *product) AddToStock(quantity int) {
	p.Quantity += quantity
}

// RemoveFromStock decreases the product's stock by the specified quantity, returns false if insufficient stock
func (p *product) RemoveFromStock(quantity int) bool {
	if p.Quantity < quantity {
		return false
	}
	p.Quantity -= quantity
	return true
}

