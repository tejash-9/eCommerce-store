package internal

import (
	"os"
	"strconv"
)

type ShoppingEngine interface {
	RegisterUser(name string, email string) (*user, error)
	GetUser(userId string) (*user, error)
	GetUserByUsername(username string) (*user, error)
	RegisterProduct(name string, description string, quantity int, sellerId string, price float64) (*product, error)
	GetProduct(productId string) (*product, error)
	AddToCart(userId string, productId string, quantity int) (map[string]int, error)
	GetCart(userId string) (map[string]int, error)
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
		interval, err := strconv.Atoi(os.Getenv(DiscountIntervalEnv))
		if err != nil {
			Logger.Sugar().Debug("Unable to read discount interval from env, using default value!")
			interval = 5 // Default discount interval if not provided
		}
		orderBook := newOrderBook()
		inventory := newInventory()
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