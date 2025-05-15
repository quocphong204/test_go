package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"order-system/config"
	"order-system/models"
)

func main() {
	config.ConnectDB()
	config.ConnectRabbitMQ()

	msgs, err := config.MQChannel.Consume(
		config.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Không thể consume từ RabbitMQ: %v", err)
	}

	fmt.Println("🐇 Đang chờ xử lý đơn hàng từ hàng đợi...")

	for msg := range msgs {
		orderIDStr := string(msg.Body)
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			log.Println("❌ Không parse được orderID:", err)
			continue
		}

		fmt.Printf("🔁 Nhận đơn hàng #%d. Đang xử lý...\n", orderID)
		time.Sleep(5 * time.Minute) // ⏳ Giả lập delay sau 5 phút (test nhanh: dùng 10s)

		var order models.Order
		if err := config.DB.First(&order, orderID).Error; err != nil {
			log.Println("❌ Không tìm thấy đơn:", err)
			continue
		}

		// Chỉ xử lý đơn còn trạng thái "pending"
		if order.Status == "pending" {
			order.Status = "processed"
			config.DB.Save(&order)
			log.Printf("✅ Đã xử lý đơn #%d (status: processed)\n", orderID)
		}
	}
}
