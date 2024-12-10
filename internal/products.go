package internal

import (
	"fmt"
)

type product struct {
	Id              string
	Name   			string
	Description    	string
	Quantity        int
	Price           float64
	SellerId        string
}

type inventory struct {
	Products   			map[string]*product
	ProductsBySeller 	map[string][]*product
}

func newProduct(id string, name string, description string, quantity int, seller_id string, price float64) *product {
	return &product {
		Id: 			id,
		Name: 			name,
		Description: 	description,
		Quantity: 		quantity,
		SellerId: 		seller_id,
		Price: 			price,
	}
}

func (s *shoppingEngine) RegisterProduct(name string, description string, quantity int, sellerId string, price float64) (*product, error) {
	for _, product := range s.Inventory.ProductsBySeller[sellerId] {
		if product.Name == name {
			return nil, fmt.Errorf("product with name already exists by the seller!")
		}
	}
	id := generateUUID()
	product := newProduct(id, name, description, quantity, sellerId, price)

	// 
	s.Inventory.ProductsBySeller[sellerId] = append(s.Inventory.ProductsBySeller[sellerId], product)
	s.Inventory.Products[id] = product
	return product, nil
}

func (p *product) IsAvailable(quantity int) bool {
	return p.Quantity >= quantity
}

func (p *product) AddtoStock(quantity int) {
	p.Quantity += quantity
}

func (p *product) RemoveFromStock(quantity int) bool {
	if p.Quantity < quantity {
		return false
	}
	p.Quantity -= quantity
	return true
}

