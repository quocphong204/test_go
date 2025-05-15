package routes

import (
	"github.com/gin-gonic/gin"
	"order-system/controllers"
	"order-system/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Public routes
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Protected routes (yêu cầu JWT)
	api := r.Group("/api")
	api.Use(middlewares.JWTAuthMiddleware()) // Gắn middleware vào nhóm này

	// Kiểm tra xác thực
	api.GET("/profile", func(c *gin.Context) {
		userID := c.GetInt("user_id")
		role := c.GetString("role")
		c.JSON(200, gin.H{
			"message": "Xác thực thành công",
			"user_id": userID,
			"role":    role,
		})
	})

	// 📦 Product CRUD routes
	api.GET("/products", controllers.GetAllProducts)         // Xem tất cả sản phẩm
	api.POST("/products", controllers.CreateProduct)         // Thêm sản phẩm (chỉ admin)
	api.PUT("/products/:id", controllers.UpdateProduct)      // Cập nhật sản phẩm (admin)
	api.DELETE("/products/:id", controllers.DeleteProduct)   // Xoá sản phẩm (admin)
	api.POST("/orders", controllers.CreateOrder)
	api.GET("/orders/me", controllers.GetMyOrders)
	api.GET("/orders", controllers.GetAllOrders)        // admin
	api.GET("/orders/:id", controllers.GetOrderByID)    // user hoặc admin
	api.GET("/orders/last", controllers.GetLastOrderFromCache)
	api.PUT("/orders/:id/pay", controllers.MarkOrderPaid)
	api.PUT("/orders/:id/process", controllers.MarkOrderProcessed)


	return r
}
