package db

import (
	"fmt"
	"os"

	rabbitmq "github.com/rabbitmq/amqp091-go"
)

// DeclareExchange sets up the exchange for publishing messages
func DeclareExchange(channel *rabbitmq.Channel) *rabbitmq.Channel {
	err := channel.ExchangeDeclare(
		"project_service.elastic_search", // exchange name
		"direct",                         // exchange type (e.g., "direct", "fanout", "topic", "headers")
		true,                             // durable
		false,                            // auto delete
		false,                            // internal
		false,                            // no wait
		nil,                              // arguments
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to RabbitMQ")
	return channel
}

// NewRabbitMQConnection creates and returns a RabbitMQ connection and channel for publishing
func NewRabbitMQConnection() *rabbitmq.Channel {
	rabbitmqServerURL := os.Getenv("RABBITMQ_SERVER_URL")

	connectRabbitMQ, err := rabbitmq.Dial(rabbitmqServerURL)
	if err != nil {
		panic(err)
	}

	// Create a channel
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		connectRabbitMQ.Close()
		panic(err)
	}

	return DeclareExchange(channelRabbitMQ)
}
