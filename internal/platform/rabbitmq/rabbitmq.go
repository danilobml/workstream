package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQClient(ctx context.Context, rabbitmqURL string, exchange string) (*RabbitMQ, error) {
	conn, err := dialWithRetry(ctx, rabbitmqURL, 10, 500*time.Millisecond)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("amqp channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		log.Printf("Failed to declare RabbitMQ exchange: %s", err)
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		_ = r.Channel.Close()
	}
	if r.Conn != nil {
		_ = r.Conn.Close()
	}
}

func dialWithRetry(
	ctx context.Context,
	url string,
	maxAttempts int,
	baseDelay time.Duration,
) (*amqp.Connection, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		conn, err := amqp.Dial(url)
		if err == nil {
			return conn, nil
		}

		lastErr = err

		delay := time.Duration(attempt) * baseDelay

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("rabbitmq connect aborted: %w", ctx.Err())
		case <-time.After(delay):
		}
	}

	return nil, fmt.Errorf("failed to connect to rabbitmq after %d attempts: %w", maxAttempts, lastErr)
}

func (r *RabbitMQ) ConsumeRabbitMQQueue(ctx context.Context, queueName string, exchange string, binding string) (<-chan amqp.Delivery, error) {
	q, err := r.Channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("Failed to declare RabbitMQ queue: %s", err)
		return nil, err
	}

	err = r.Channel.QueueBind(
		q.Name,   // queue name
		binding,  // routing key
		exchange, // exchange
		false,    // no wait
		nil,      // args
	)
	if err != nil {
		log.Printf("Failed to bind queue: %s", err)
		return nil, err
	}

	msgs, err := r.Channel.ConsumeWithContext(
		ctx,    // context
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Printf("Failed to register RabbitMQ consumer: %s", err)
		return nil, err
	}

	return msgs, nil
}
