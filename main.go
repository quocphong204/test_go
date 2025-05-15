package main

import (
	"order-system/config"
	"order-system/routes"
)

func main() {
	config.ConnectDB()
	config.ConnectRedis()
	config.ConnectRabbitMQ()
	r := routes.SetupRouter()
	r.Run(":8080")
}
