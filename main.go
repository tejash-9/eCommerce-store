package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ecommerce-store/internal"
	"github.com/ecommerce-store/routes"
)

func main() {
	// Get App instance
	shoppingApp := internal.GetAppInstance()

	// Initialize gin Router
	route := gin.Default()

	routes.RegisterRoutes(route, shoppingApp)
	
	route.Run(":8085")

}