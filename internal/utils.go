package internal

import (
	"sync"
	"math/rand"
    "time"
    "strings"
    "encoding/base64"
    "github.com/google/uuid"
)

var (
	shoppingApp 	*shoppingEngine
	instance        sync.Once
)

// nTH order
const OrderThresholdEnv = "ORDER_THRESHOLD"

// GenerateUUID generates a new UUID
func generateUUID() string {
    // Create a new UUID
    newUUID := uuid.New()

    // Encode the UUID to Base64 (URL-safe, without padding)
	encoded := base64.RawURLEncoding.EncodeToString(newUUID[:]) // RawURLEncoding for URL-safe encoding
	return encoded
}