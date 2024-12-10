package routes

import (
	"github.com/ecommerce-store/internal"
	"github.com/gin-gonic/gin"
)


func RegisterRoutes(router *gin.Engine, app internal.ShoppingEngine) {
	admin := router.Group("/admin")
	registerAdminRoutes(admin, app)

	auth := router.Group("/auth")
	registerAuthRoutes(auth, app)

	users := router.Group("/users")
	registerUserRoutes(users, app)

	orders := router.Group("orders")
	registerOrderRoutes(orders, app)

	products := router.Group("/products")
	registerProductRoutes(products, app)
}

func registerAdminRoutes(rg *gin.RouterGroup, app internal.ShoppingEngine) {
	rg.GET("/total_purchases", func(c *gin.Context) {
		// Get the total amount for the purchase
		amount := app.OrderHistory().TotalPurchaseAmount()

		// Check if amount is valid, in case the order book is empty or unavailable
		if amount < 0.0 {
			c.JSON(500, gin.H{
				"status":  "failed",
				"message": "Error fetching total purchase amount",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  	"success",
			"message": 	"Total purchase amount retrieved successfully",
			"body":     gin.H{
				"purchase_amount": amount,
			},
		})
	})

	rg.GET("/coupons", func(c *gin.Context) {
		// List all available discount codes
		coupons := app.OrderHistory().ListDiscountCoupons()

		// Check if no codes found or something went wrong
		if coupons == nil || len(coupons) == 0 {
			c.JSON(404, gin.H{
				"status":  "failed",
				"message": "No discount coupons found",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Discount coupons retrieved successfully",
			"body":     gin.H{
				"discount_coupons": coupons,
			},
		})
	})

	rg.GET("/items", func(c *gin.Context) {
		items := app.OrderHistory().TotalSoldItems()

		if items < 0 {
			c.JSON(500, gin.H{
				"status":  "failed",
				"message": "Error fetching total sold Items",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":         "success",
			"message":        "Total items sold retrieved successfully",
			"body":           gin.H{
				"items_sold": items,
			},
		})
	})

	rg.GET("/total_discount", func(c *gin.Context) {
		// Get the total discount applied
		totalDiscount := app.OrderHistory().TotalDiscountAmount()

		// Check if discount is valid
		if totalDiscount < 0 {
			c.JSON(500, gin.H{
				"status":  "failed",
				"message": "Error fetching total discount amount",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":         "success",
			"message":        "Total discount amount retrieved successfully",
			"body":           gin.H{
				"total_discount": totalDiscount,
			},
		})
	})

}

func registerAuthRoutes(rg *gin.RouterGroup, app internal.ShoppingEngine) {
	rg.POST("/login", func(c *gin.Context) {
		var request struct {
			Username string `json:"username"`
		}
	
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "Invalid request format",
			})
			return
		}
	
		// Check if the username and password match
		user, err := app.GetUserByUsername(request.Username)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Login successful",
			"body":    gin.H{"user": user},
		})
	})

	rg.POST("/register", func(c *gin.Context) {
		var request struct {
			Email 	 string `json:"email"`
			Name     string `json:"name"`
		}
	
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "Invalid request format",
			})
			return
		}
	
		if request.Email == "" || request.Name == "" {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "name or email is empty",
			})
			return
		}

		user, err := app.RegisterUser(request.Name, request.Email)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"status":  "success",
			"message": "User added successfully",
			"body":    gin.H{"user": user},
		})
	})
}

func registerUserRoutes(rg *gin.RouterGroup, app internal.ShoppingEngine) {

	user := rg.Group("/:user_id")
	user.GET("/", func(c *gin.Context) {
		// Parse user id from the URL parameters (e.g., /:user_id/cart)
		userId := c.Param("user_id")
		if userId == "" {
			// If userId is empty, return a bad request error
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "Invalid user Id",
			})
		}

		user, err := app.GetUser(userId)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"status":  "success",
			"message": "User retrieved successfully",
			"body":    gin.H{
				"user": user,
			},
		})
	})

	user.GET("/coupon", func(c *gin.Context) {
		// Parse user id from the URL parameters (e.g., /:user_id/cart)
		userId := c.Param("user_id")
		if userId == "" {
			// If userId is empty, return a bad request error
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "Invalid user Id",
			})
		}

		coupon, err := app.GetDiscountCoupon(userId)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}

		c.JSON(201, gin.H{
			"status":  "success",
			"message": "Discount coupon retrieved successfully",
			"body":    gin.H{
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
				"status":  "failed",
				"message": "Invalid user Id",
			})
		}

		var request struct {
			ProductId  string  `json:"product_id"`
			Quantity   int     `json:"quantity"`
		}
	
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "Invalid request format",
			})
			return
		}

		cart, err := app.AddToCart(userId, request.ProductId, request.Quantity)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Product added to cart successfully",
			"body":    gin.H{"cart": cart},
		})
	})
}

func registerProductRoutes(rg *gin.RouterGroup, app internal.ShoppingEngine) {

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
				"status":  "failed",
				"message": "Invalid request format",
			})
			return
		}
	
		// Ensure userId is provided
		if request.UserId == "" || request.Name == "" || request.Price < 0.0 || request.Quantity < 0 {
			c.JSON(400, gin.H{"error": "User ID is required"})
			return
		}
	
		// Add the product (call to backend logic)
		product, err := app.RegisterProduct(request.Name, request.Description, request.Quantity, request.UserId, request.Price)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}
	
		c.JSON(201, gin.H{
			"status":  "success",
			"message": "Product added successfully",
			"body":    gin.H{"product": product},
		})
	})
	
	rg.GET("/:product_id", func(c *gin.Context) {
		// Parse user id from the URL parameters (e.g., /:user_id/cart)
		productId := c.Param("product_id")
		if productId == "" {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "Invalid Product Id",
			})
		}
	
		product, err := app.GetProduct(productId)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}
	
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Product retrieved successfully",
			"body":    gin.H{
				"product": product,
			},
		})
	})
}

func registerOrderRoutes(rg *gin.RouterGroup, app internal.ShoppingEngine) {

	rg.POST("/checkout", func(c *gin.Context) {
		// Parse coupon code from the request body (assuming it's a JSON request)
		var request struct {
			UserId     string `json:"user_id"`
			CouponCode string `json:"coupon_code"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": "Invalid request format",
			})
			return
		}
		
		// Call the Checkout function (assumed to be a method of an object 's')
		order, err := app.Checkout(request.UserId, request.CouponCode)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "failed",
				"message": err.Error(),
			})
			return
		}
	
		// Return the order as a successful response
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Order placed successfully",
			"body":    gin.H{"order": order},
		})
	})
}
