package queues

import (
	"encoding/json"
	"fmt"

	rabbitmq "github.com/rabbitmq/amqp091-go"
)

// Message represents the structure of the message to be sent to RabbitMQ
type Message struct {
	Operation string      `json:"operation"` // e.g., "create", "update", "delete"
	Data      interface{} `json:"data"`      // The actual data related to the operation
}

func PublishMessageFromRabbitMQToElasticSearch(channel *rabbitmq.Channel, operation string, data interface{}) error {
	message := Message{
		Operation: operation,
		Data:      data,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"project_service.search",      // exchange
		"project_service.events.crud", // routing key
		false,                         // mandatory
		false,                         // immediate
		rabbitmq.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	} else {
		fmt.Println("Message published successfully.")
	}

	return err
}
