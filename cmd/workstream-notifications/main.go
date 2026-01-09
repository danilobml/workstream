package main

import (
	"context"
	"log"
	"os"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/platform/rabbitmq"
	"github.com/danilobml/workstream/internal/workstream-notifications/readiness"
	"github.com/danilobml/workstream/internal/workstream-notifications/services"
)

const (
	serviceName     = "workstream-notifications"
	httpPortName    = "NOTIFICATIONS_HTTP_PORT"
	rabbitmqUrlName = "RABBITMQ_URL"
)

func main() {
	rabbitmqUrl := os.Getenv(rabbitmqUrlName)
	if rabbitmqUrl == "" {
		log.Fatal("unable to read RABBITMQ_URL from env")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rabbitClient, err := rabbitmq.NewRabbitMQClient(ctx, rabbitmqUrl, rabbitmq.Exchange)
	if err != nil {
		log.Fatal("workstream-notifications - failed to connect to RabbitMQ", err)
	}
	defer rabbitClient.Close()

	rabbitService := services.NewRabbitConsumerService(rabbitClient)

	go func() {
		if err := rabbitService.Consume(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	if err := http.StartServer(
		serviceName,
		httpPortName,
		nil,
		readiness.IsReady,
	); err != nil {
		log.Fatal(err)
	}
}
