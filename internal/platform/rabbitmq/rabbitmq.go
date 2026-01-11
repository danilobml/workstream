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

func (r *RabbitMQ) ConsumeRabbitMQQueue(ctx context.Context, queueName string) (<-chan amqp.Delivery, error) {
	msgs, err := r.Channel.ConsumeWithContext(
		ctx,
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("consume %s: %w", queueName, err)
	}
	return msgs, nil
}

func (r *RabbitMQ) DeclareQueues(mainQueueName, exchange, binding string) error {
	dlx := fmt.Sprintf("%s.dlx", mainQueueName)
	dlq := fmt.Sprintf("%s.dlq", mainQueueName)
	
	// Declare DLX and DLQ
	if err := r.Channel.ExchangeDeclare(dlx, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare dlx %s: %w", dlx, err)
	}

	if _, err := r.Channel.QueueDeclare(dlq, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declare dlq %s: %w", dlq, err)
	}

	if err := r.Channel.QueueBind(dlq, dlq, dlx, false, nil); err != nil {
		return fmt.Errorf("bind dlq %s to dlx %s: %w", dlq, dlx, err)
	}

	// Declare main queue with DLX args
	args := amqp.Table{
		"x-dead-letter-exchange":    dlx,
		"x-dead-letter-routing-key": dlq,
	}

	q, err := r.Channel.QueueDeclare(mainQueueName, true, false, false, false, args)
	if err != nil {
		return fmt.Errorf("declare main queue %s: %w", mainQueueName, err)
	}

	if err := r.Channel.QueueBind(q.Name, binding, exchange, false, nil); err != nil {
		return fmt.Errorf("bind main queue %s: %w", mainQueueName, err)
	}

	log.Printf("RabbitMQ main and dead letter queues for %s successfully declared and bound!", mainQueueName)

	return nil
}
