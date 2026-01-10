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

type MessageConsumerService interface {
	Consume()
	ProcessEvent(ctx context.Context, event models.Event) error
}

type RabbitMessageConsumerService struct {
	client              *rabbitmq.RabbitMQ
	eventsProcessor EventsProcessor
}

func NewRabbitMessageConsumerService(client *rabbitmq.RabbitMQ, eventsProcessor EventsProcessor) *RabbitMessageConsumerService {
	return &RabbitMessageConsumerService{
		client:              client,
		eventsProcessor: eventsProcessor,
	}
}

func (rs *RabbitMessageConsumerService) Consume(ctx context.Context) error {
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

func (rs *RabbitMessageConsumerService) ProcessEvent(ctx context.Context, event models.Event) error {

	log.Printf("Received a message = %v\n", event)

	err := rs.eventsProcessor.SaveNewEvent(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
