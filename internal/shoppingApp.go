package internal

import (
	"os"
	"log"
	"fmt"
	"sync"
	"strconv"
)

type ShoppingEngine interface {
	NewUser(name string, email string) (*user, error)
	GetUserInfoByUsername(username string) (*user, error)
}

type shoppingEngine struct {
	Users   			map[string]*user
	UserMap         	map[string]string
	OrderThreshold  	int
	CouponMutex			*sync.Mutex
}

func GetShoppingInstance() ShoppingEngine {
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
			OrderThreshold: threshold,
			CouponMutex: 	&sync.Mutex{},
		}
	})
	return shoppingApp
}