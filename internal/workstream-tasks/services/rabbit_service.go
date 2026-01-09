package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/danilobml/workstream/internal/platform/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventsService interface {
	Publish(ctx context.Context, event models.Event) error
}

type RabbitProducerService struct {
	client *rabbitmq.RabbitMQ
}

func NewRabbitProducerService(client *rabbitmq.RabbitMQ) *RabbitProducerService {
	return &RabbitProducerService{
		client: client,
	}
}

func (rs *RabbitProducerService) Publish(ctx context.Context, event models.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	routingKey := event.EventType

	err = rs.client.Channel.PublishWithContext(
		ctx,               // context
		rabbitmq.Exchange, // exchange
		routingKey,        // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		})

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	log.Printf("Message published to exchange: %s with key: %s", rabbitmq.Exchange, routingKey)

	return nil
}
