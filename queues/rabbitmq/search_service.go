package queues

import (
	"encoding/json"
	"fmt"
	"log"

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
		log.Println("Error marshalling message:", err)
		return err
	}
	fmt.Printf("%+v\n", message)
	err = channel.Publish(
		"project_service.elastic_search", // exchange
		"elastic_search.crud",            // routing key
		false,                            // mandatory
		false,                            // immediate
		rabbitmq.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}
