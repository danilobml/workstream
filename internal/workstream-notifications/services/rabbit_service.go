package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/danilobml/workstream/internal/platform/errs"
	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/danilobml/workstream/internal/platform/rabbitmq"
)

type EventConsumerService interface {
	Consume()
	ProcessEvent(ctx context.Context, event models.Event) error
}

type RabbitConsumerService struct {
	client              *rabbitmq.RabbitMQ
	notificationService NotificationService
}

func NewRabbitConsumerService(client *rabbitmq.RabbitMQ, notificationService NotificationService) *RabbitConsumerService {
	return &RabbitConsumerService{
		client:              client,
		notificationService: notificationService,
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

		err := rs.ProcessEvent(ctx, event)
		if errors.Is(err, errs.ErrAlreadyProcessed) {
			log.Printf("process failed (skip): %v", err)
			if ackErr := d.Ack(false); ackErr != nil {
				log.Printf("ack failed: %v", ackErr)
			}
			continue
		}
		if err != nil {
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

func (rs *RabbitConsumerService) ProcessEvent(ctx context.Context, event models.Event) error {

	log.Printf("Received a message = %v\n", event)

	err := rs.notificationService.CreateNewNotification(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
