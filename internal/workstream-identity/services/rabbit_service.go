package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/danilobml/workstream/internal/platform/models"
	"github.com/danilobml/workstream/internal/platform/rabbitmq"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventsService interface {
	Publish(ctx context.Context, event models.Event) error
	SendMailMessage(ctx context.Context, user models.User, subject, body string) error
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
		rabbitmq.NotificationsExchange, // exchange
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

	log.Printf("Message published to exchange: %s with key: %s", rabbitmq.NotificationsExchange, routingKey)

	return nil
}


func (ns *RabbitProducerService) SendMailMessage(ctx context.Context, user models.User, subject, body string) error {
	mailInput := models.MailInput{
		To:      []string{user.Email},
		Subject: subject,
		Body:    body,
	}
	payloadBytes, err := json.Marshal(mailInput)
	if err != nil {
		return status.Error(codes.Internal, "failed to marshal mail payload")
	}

	event := models.Event{
		EventID:    uuid.NewString(),
		EventType:  "workstream.mail.sent.v1",
		OccurredAt: time.Now(),
		Producer:   "workstream-notifications",
		TraceID:    uuid.NewString(),
		Payload:    json.RawMessage(payloadBytes),
	}

	err = ns.Publish(ctx, event)
	if err != nil {
		return status.Error(codes.Internal, "failed to send event")
	}

	return nil
}
