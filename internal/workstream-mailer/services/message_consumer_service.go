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
	Consume(ctx context.Context) error
	ProcessEvent(ctx context.Context, mailInput models.MailInput) error
}

type RabbitMessageConsumerService struct {
	client      *rabbitmq.RabbitMQ
	mailService Mailer
}

func NewRabbitMessageConsumerService(client *rabbitmq.RabbitMQ, mailer Mailer) *RabbitMessageConsumerService {
	return &RabbitMessageConsumerService{
		client:      client,
		mailService: mailer,
	}
}

func (rs *RabbitMessageConsumerService) Consume(ctx context.Context) error {
	msgs, err := rs.client.ConsumeRabbitMQQueue(ctx, rabbitmq.MailerQueue)
	if err != nil {
		return err
	}

	log.Printf("[*] Waiting for messages...")

	for d := range msgs {
		var evt models.Event
		if err := json.Unmarshal(d.Body, &evt); err != nil {
			log.Printf("invalid event payload (DLQ): %v body=%s", err, string(d.Body))
			_ = d.Nack(false, false)
			continue
		}

		if evt.EventID == "" || evt.EventType == "" {
			log.Printf(
				"invalid event - missing required fields (DLQ): event_id=%q event_type=%q",
				evt.EventID,
				evt.EventType,
			)
			_ = d.Nack(false, false)
			continue
		}

		var mailInput models.MailInput
		if err := json.Unmarshal(evt.Payload, &mailInput); err != nil {
			log.Printf(
				"invalid mail payload (DLQ): %v payload=%s",
				err,
				string(evt.Payload),
			)
			_ = d.Nack(false, false)
			continue
		}

		log.Printf(
			"decoded mail payload: to=%v subject=%q",
			mailInput.To,
			mailInput.Subject,
		)

		err := rs.ProcessEvent(ctx, mailInput)

		if errors.Is(err, errs.ErrAlreadyProcessed) {
			log.Printf(
				"event skipped (already processed): event_id=%s trace_id=%s",
				evt.EventID,
				evt.TraceID,
			)
			_ = d.Ack(false)
			continue
		}

		if errors.Is(err, errs.ErrInProgress) {
			log.Printf(
				"still in progress (requeue): event_id=%s trace_id=%s err=%v",
				evt.EventID,
				evt.TraceID,
				err,
			)
			_ = d.Nack(false, true)
			continue
		}

		if errors.Is(err, errs.ErrInvalidEvent) {
			log.Printf(
				"invalid event (DLQ): event_id=%s trace_id=%s err=%v",
				evt.EventID,
				evt.TraceID,
				err,
			)
			_ = d.Nack(false, false)
			continue
		}

		if err != nil {
			log.Printf(
				"process failed (requeue): event_id=%s trace_id=%s err=%v",
				evt.EventID,
				evt.TraceID,
				err,
			)
			_ = d.Nack(false, true)
			continue
		}

		if err := d.Ack(false); err != nil {
			log.Printf("ack failed: %v", err)
		}
	}

	return nil
}

func (rs *RabbitMessageConsumerService) ProcessEvent(ctx context.Context,mailInput models.MailInput) error {
	return rs.mailService.SendMail(mailInput)
}
