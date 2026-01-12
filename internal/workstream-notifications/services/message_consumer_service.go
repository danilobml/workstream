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
	client          *rabbitmq.RabbitMQ
	eventsProcessor EventsProcessor
}

func NewRabbitMessageConsumerService(client *rabbitmq.RabbitMQ, eventsProcessor EventsProcessor) *RabbitMessageConsumerService {
	return &RabbitMessageConsumerService{
		client:          client,
		eventsProcessor: eventsProcessor,
	}
}

func (rs *RabbitMessageConsumerService) Consume(ctx context.Context) error {
	msgs, err := rs.client.ConsumeRabbitMQQueue(ctx, rabbitmq.Queue)
	if err != nil {
		return err
	}

	log.Printf("[*] Waiting for messages...")

	for d := range msgs {
		var event models.Event

		if err := json.Unmarshal(d.Body, &event); err != nil {
			log.Printf("invalid event payload (DLQ): %v", err)
			if nackErr := d.Nack(false, false); nackErr != nil {
				log.Printf("nack failed: %v", nackErr)
			}
			continue
		}

		if event.EventID == "" || event.EventType == "" {
			log.Printf("invalid event - missing required fields (DLQ): event_id=%q event_type=%q", event.EventID, event.EventType)
			if nackErr := d.Nack(false, false); nackErr != nil {
				log.Printf("nack failed: %v", nackErr)
			}
			continue
		}

		err := rs.ProcessEvent(ctx, event)

		if errors.Is(err, errs.ErrAlreadyProcessed) {
			log.Printf("event skipped (already processed): event_id=%s trace_id=%s", event.EventID, event.TraceID)
			if ackErr := d.Ack(false); ackErr != nil {
				log.Printf("ack failed: %v", ackErr)
			}
			continue
		}

		if errors.Is(err, errs.ErrInProgress) {
			log.Printf("still in progress (requeue): event_id=%s trace_id=%s err=%v", event.EventID, event.TraceID, err)
			if nackErr := d.Nack(false, true); nackErr != nil {
				log.Printf("nack failed: %v", nackErr)
			}
			continue
		}

		if errors.Is(err, errs.ErrInvalidEvent) {
			log.Printf("invalid event (DLQ): event_id=%s trace_id=%s err=%v", event.EventID, event.TraceID, err)
			if nackErr := d.Nack(false, false); nackErr != nil {
				log.Printf("nack failed: %v", nackErr)
			}
			continue
		}

		if err != nil {
			log.Printf("process failed (requeue): event_id=%s trace_id=%s err=%v", event.EventID, event.TraceID, err)
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

	err := rs.eventsProcessor.ProcessEvent(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
