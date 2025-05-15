package config

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var MQConn *amqp.Connection
var MQChannel *amqp.Channel
var QueueName = "order_queue"
var MQCtx = context.Background()
var MQTimeout = 5 * time.Second

func ConnectRabbitMQ() {
	var err error
	MQConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("❌ Không thể kết nối RabbitMQ:", err)
	}

	MQChannel, err = MQConn.Channel()
	if err != nil {
		log.Fatal("❌ Không thể tạo channel:", err)
	}

	_, err = MQChannel.QueueDeclare(
		QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("❌ Không thể tạo queue:", err)
	}

	fmt.Println("✅ Kết nối RabbitMQ thành công")
}
