package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/danilobml/workstream/internal/platform/rabbitmq"
)

type EventConsumerService interface {
	Consume()
	ProcessEvent(event models.Event) error
}

type RabbitConsumerService struct {
	client    *rabbitmq.RabbitMQ
}

func NewRabbitConsumerService(client *rabbitmq.RabbitMQ) *RabbitConsumerService {
	return &RabbitConsumerService{
		client: client,
	}
}

func (rs *RabbitConsumerService) Consume(ctx context.Context) error {
	msgs, err := rs.client.ConsumeRabbitMQQueue(ctx, rabbitmq.Queue, rabbitmq.Exchange, rabbitmq.Binding)
	if err != nil {
		return err
	}

	log.Printf("[*] Waiting for messages...")

	for d := range msgs {
		var event models.Event

		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("invalid event payload: %v", err)
			if ackErr := d.Ack(false); ackErr != nil {
				log.Printf("ack failed: %v", ackErr)
			}
			continue
		}

		if err := rs.ProcessEvent(event); err != nil {
			log.Printf("process failed (requeue): %v", err)
			if nackErr := d.Nack(false, true); nackErr != nil {
				log.Printf("nack failed: %v", nackErr)
			}
			continue
		}

		if ackErr := d.Ack(false); ackErr != nil {
			log.Printf("ack failed: %v", ackErr)
		}
	}

	return nil
}


func (rs *RabbitConsumerService) ProcessEvent(event models.Event) error {
	// TODO - implement functionality:
	log.Printf("Received a message = %v\n", event)

	return nil
}
