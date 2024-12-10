package internal

// order represents an order placed by a user
type order struct {
	Id              string           	`json:"id"`              	// Unique order ID
	UserId          string          	`json:"user_id"`         	// User ID who placed the order
	OrderCart       map[string]int  	`json:"order_cart"`      	// Cart with product IDs and quantities
	CartTotal     	float64             `json:"amount"`           	// Total cart value before discount
	Discount        float64            	`json:"discount"`        	// Discount applied on the order
	DiscountCoupon  string            	`json:"discount_coupon"` 	// Applied coupon code
	AmountToPay     float64          	`json:"amount_to_pay"`    	// Final amount after discount
}

// newOrder creates a new order instance
func newOrder(id string, userId string, cart map[string]int, amount float64, coupon string, discount float64, finalAmount float64) *order {
	return &order{
		Id:             id,               // Set unique order ID
		UserId:         userId,           // Set user ID
		OrderCart:      cart,             // Set order cart with product quantities
		CartTotal:      amount,           // Set total cart amount before discount
		DiscountCoupon: coupon,           // Set applied coupon code
		Discount:       discount,         // Set discount applied on the order
		AmountToPay:    finalAmount,      // Set final amount after discount
	}
}