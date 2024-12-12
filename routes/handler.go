package routes

import (
	"net/mail"
	"github.com/ecommerce-store/internal"
	"github.com/gin-gonic/gin"
)


func RegisterRoutes(router *gin.Engine, svc internal.ShoppingEngine) {
	admin := router.Group("/admin")
	registerAdminRoutes(admin, svc)

	auth := router.Group("/auth")
	registerAuthRoutes(auth, svc)

	users := router.Group("/users")
	registerUserRoutes(users, svc)

	orders := router.Group("/orders")
	registerOrderRoutes(orders, svc)

	products := router.Group("/products")
	registerProductRoutes(products, svc)
}

func registerAdminRoutes(rg *gin.RouterGroup, svc internal.ShoppingEngine) {
	rg.GET("/analytics", func(c *gin.Context) {
		// Get the total amount for the purchase
		items, amount, discount, coupons := svc.OrderHistory().GetAnalytics()
		c.JSON(200, gin.H{
			"status":  	"success",
			"message": 	"Platform analytics retrieved successfully",
			"data": gin.H{
				"total_items_sold": items,
				"total_purchase_amount": amount,
				"total_discount": discount,
				"applied_coupons": coupons,
			},
		})
	})
}

func registerAuthRoutes(rg *gin.RouterGroup, svc internal.ShoppingEngine) {
	rg.POST("/login", func(c *gin.Context) {
		var request struct {
			Username string `json:"username"`
		}
	
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Invalid request",
			})
			return
		}
	
		// Check if the username and password match
		user, err := svc.GetUserByUsername(request.Username)
		if err != nil {
			c.JSON(404, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Login successful",
			"data":    gin.H{
				"user": gin.H{
					"id": user.Id,
					"name": user.Name,
					"email": user.Email,
				},
			},
		})
	})

	rg.POST("/register", func(c *gin.Context) {
		var request struct {
			Email 	string `json:"email"`
			Name    string `json:"name"`
		}
	
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Invalid request",
			})
			return
		}
	
		if request.Email == "" || request.Name == "" {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Name and email are required",
			})
			return
		}

		// Validate email format
		_, err := mail.ParseAddress(request.Email)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Invalid email format",
			})
			return
		}

		user, err := svc.RegisterUser(request.Name, request.Email)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"status":  "success",
			"message": "User registered successfully",
			"data":    gin.H{
				"user": gin.H{
					"id": user.Id,
					"name": user.Name,
					"email": user.Email,
				},
			},
		})
	})
}

func registerUserRoutes(rg *gin.RouterGroup, svc internal.ShoppingEngine) {

	user := rg.Group("/:user_id")

	user.GET("/coupon", func(c *gin.Context) {
		// Parse user id from the URL parameters (e.g., /:user_id/cart)
		userId := c.Param("user_id")
		if userId == "" {
			// If userId is empty, return a bad request error
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "User ID cannot be empty",
			})
			return
		}

		// Get discount coupon for the user
		coupon, err := svc.GetDiscountCoupon(userId)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Discount coupon retrieved successfully",
			"data":    gin.H{
				"coupon_code": coupon,
			},
		})
	})

	user.POST("/cart", func(c *gin.Context) {
		// Parse user id from the URL parameters (e.g., /:user_id/cart)
		userId := c.Param("user_id")
		if userId == "" {
			// If userId is empty, return a bad request error
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "User ID cannot be empty",
			})
			return
		}

		var cartItem struct {
			ProductId  string  `json:"product_id"`
			Quantity   int     `json:"quantity"`
		}
	
		if err := c.ShouldBindJSON(&cartItem); err != nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Invalid request",
			})
			return
		}

		cartMap, err := svc.AddToCart(userId, cartItem.ProductId, cartItem.Quantity)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		var updatedCart []gin.H
		for key, value := range cartMap {
			updatedCart = append(updatedCart, gin.H{
				"product_id": 	key,
				"quantity": 	value,
			})
		}
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Product added to cart successfully",
			"data":    gin.H{
				"user_id": 	userId,
				"cart": 	updatedCart,
			},
		})
	})

	user.GET("/cart", func(c *gin.Context) {
		// Parse user id from the URL parameters (e.g., /:user_id/cart)
		userId := c.Param("user_id")
		if userId == "" {
			// If userId is empty, return a bad request error
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "User ID cannot be empty",
			})
			return
		}

		cartMap, err := svc.GetCart(userId)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		var updatedCart []gin.H
		for key, value := range cartMap {
			updatedCart = append(updatedCart, gin.H{
				"product_id": key,
				"quantity": value,
			})
		}
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "User cart retrieved successfully",
			"data":    gin.H{
				"user_id": 	userId,
				"cart": 	updatedCart,
			},
		})
	})
}

func registerProductRoutes(rg *gin.RouterGroup, svc internal.ShoppingEngine) {

	rg.POST("/", func(c *gin.Context) {
		var request struct {
			UserId      string  `json:"user_id"`
			Name        string  `json:"name"`
			Description string  `json:"description"`
			Price       float64 `json:"price"`
			Quantity    int   	`json:"quantity"`
		}
	
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Invalid request format",
			})
			return
		}
	
		// validate the request
		if request.UserId == "" || request.Name == "" || request.Price < 0.0 || request.Quantity < 0 {
			c.JSON(400, gin.H{
				"status": "error",
				"message": "Invalid request data",
			})
			return
		}
	
		// Add the product (call to backend logic)
		product, err := svc.RegisterProduct(request.Name, request.Description, request.Quantity, request.UserId, request.Price)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	
		c.JSON(201, gin.H{
			"status":  "success",
			"message": "Product added successfully",
			"data":    gin.H{
				"product": product,
			},
		})
	})
	
	rg.GET("/:product_id", func(c *gin.Context) {
		// Parse user id from the URL parameters (e.g., /:user_id/cart)
		productId := c.Param("product_id")
		if productId == "" {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Product ID cannot be empty",
			})
			return
		}
		
		product, err := svc.GetProduct(productId)
		if err != nil {
			c.JSON(404, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Product retrieved successfully",
			"data":    gin.H{
				"product": product,
			},
		})
	})
}

func registerOrderRoutes(rg *gin.RouterGroup, svc internal.ShoppingEngine) {

	rg.POST("/checkout", func(c *gin.Context) {
		// Parse coupon code from the request body (assuming it's a JSON request)
		var request struct {
			UserId     string `json:"user_id"`
			CouponCode string `json:"coupon_code"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "Invalid request",
			})
			return
		}
		if request.UserId == "" {
			c.JSON(400, gin.H{
				"status":  "error",
				"message": "User ID cannot be empty",
			})
			return
		}
		
		// Call the Checkout function
		order, err := svc.Checkout(request.UserId, request.CouponCode)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	
		// Return the order as a successful response
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Order placed successfully",
			"data":    gin.H{
				"order": order,
			},
		})
	})
}
