package internal

import (
	"sync"
	"math/rand"
    "time"
    "strings"
    "encoding/base64"
    "github.com/google/uuid"
    "github.com/ecommerce-store/utilities"
)

var (
	shoppingApp 	*shoppingEngine
	instance        sync.Once
    Logger = utilities.Logger.Session("dev", "github.com/ecommerce-store")
)

// nTH order (value of n)
const DiscountIntervalEnv = "DISCOUNT_INTERVAL"

// GenerateUUID generates a new UUID
func generateUUID() string {
    // Create a new UUID
    newUUID := uuid.New()

    // Encode the UUID to Base64 (URL-safe, without padding)
	encoded := base64.RawURLEncoding.EncodeToString(newUUID[:]) // RawURLEncoding for URL-safe encoding
	return encoded
}

// GenerateCouponCode generates a random alphanumeric coupon code
func generateCouponCode(length int) string {
    // Define the characters that can be used in the coupon code
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

    // Initialize random number generator
    rand.Seed(time.Now().UnixNano())

    // Generate a random string of the specified length
    var builder strings.Builder
    for i := 0; i < length; i++ {
        builder.WriteByte(chars[rand.Intn(len(chars))])
    }

    return builder.String()
}