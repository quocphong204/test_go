package controllers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"order-system/config"
	"order-system/models"
	"order-system/utils"
	"github.com/gin-gonic/gin"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

func CreateOrder(c *gin.Context) {
	userID := c.GetInt("user_id")
	fmt.Println("✅ userID:", userID)

	var request struct {
		Items []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
		return
	}

	var orderItems []models.OrderItem
	var total float64

	for _, item := range request.Items {
		var product models.Product
		if err := config.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Sản phẩm ID %d không tồn tại", item.ProductID)})
			return
		}

		unitPrice := product.Price
		totalPrice := product.Price * float64(item.Quantity)

		orderItems = append(orderItems, models.OrderItem{
			ProductID:  product.ID,
			Quantity:   item.Quantity,
			UnitPrice:  unitPrice,
			TotalPrice: totalPrice,
		})

		total += totalPrice
	}

	order := models.Order{
	UserID:        uint(userID),
	Items:         orderItems,
	TotalPrice:    total,
	Status:        "pending",         // nếu bạn có xử lý trạng thái đơn
	PaymentStatus: "unpaid",          // ✅ thêm dòng này
	}


	if err := config.DB.Create(&order).Error; err != nil {
	fmt.Println("❌ DB error:", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo đơn hàng"})
	return
	}
	// 🔔 Gửi email xác nhận đơn hàng
var user struct {
	Email string
	Name  string
}

if err := config.DB.Model(&models.User{}).
	Select("email", "name").
	Where("id = ?", userID).
	Scan(&user).Error; err != nil {
	fmt.Println("❌ Không thể lấy thông tin người dùng để gửi mail:", err)
}


	// Gửi mail
	_ = utils.SendOrderConfirmationEmail(user.Email, user.Name, order.ID, order.TotalPrice)

	// 🔁 Cache đơn hàng vào Redis
	orderJSON, _ := json.Marshal(order)
	key := fmt.Sprintf("last_order_user_%d", userID)
	config.RedisClient.Set(config.Ctx, key, orderJSON, 10*time.Minute)

	// 🐇 Gửi đơn hàng vào hàng đợi RabbitMQ để xử lý sau 5 phút
message := fmt.Sprintf("%d", order.ID)
err := config.MQChannel.Publish(
	"",                   // default exchange
	config.QueueName,     // queue name
	false, false,
	amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	},
)
if err != nil {
	fmt.Println("❌ Lỗi khi gửi đơn vào RabbitMQ:", err)
}

	
	c.JSON(http.StatusCreated, order)
	
}

func GetMyOrders(c *gin.Context) {
	userID := c.GetInt("user_id")

	var orders []models.Order
	err := config.DB.Preload("Items.Product").Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy đơn hàng"})
		return
	}

		
	c.JSON(http.StatusOK, orders)
}

func GetAllOrders(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Chỉ admin được phép xem tất cả đơn hàng"})
		return
	}

	var orders []models.Order
	if err := config.DB.Preload("User").Preload("Items.Product").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể lấy danh sách đơn hàng"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetInt("user_id")
	role := c.GetString("role")

	var order models.Order
	if err := config.DB.Preload("User").Preload("Items.Product").First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đơn hàng"})
		return
	}

	// Chỉ cho phép admin hoặc chính user xem đơn của mình
	if role != "admin" && int(order.UserID) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bạn không có quyền truy cập đơn hàng này"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func GetLastOrderFromCache(c *gin.Context) {
	userID := c.GetInt("user_id")
	key := fmt.Sprintf("last_order_user_%d", userID)

	data, err := config.RedisClient.Get(config.Ctx, key).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đơn hàng trong cache"})
		return
	}

	var order models.Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi parse dữ liệu Redis"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func MarkOrderPaid(c *gin.Context) {
	orderID := c.Param("id")

	var order models.Order
	if err := config.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đơn hàng"})
		return
	}

	order.PaymentStatus = "paid"

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật trạng thái thanh toán"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đã thanh toán thành công"})
}

func MarkOrderProcessed(c *gin.Context) {
	orderID := c.Param("id")

	var order models.Order
	if err := config.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy đơn hàng"})
		return
	}

	if order.Status == "processed" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Đơn hàng đã được xử lý trước đó"})
		return
	}

	order.Status = "processed"

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể cập nhật trạng thái đơn hàng"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Đơn hàng đã được xử lý"})
}

