package main

import (
	"os"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/ecommerce-store/internal"
	"github.com/ecommerce-store/routes"
)

func main() {
	// Get App instance
	instance := internal.GetAppInstance()

	// Initialize gin Router
	route := gin.Default()

	// Register routes
	routes.RegisterRoutes(route, instance)

	// Loading port from env
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("Port not found, Running on default 8080...")
		port = "8080"
	}
	
	// Listen and serve on the specified port
	err := route.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}