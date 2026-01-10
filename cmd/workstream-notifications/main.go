package main

import (
	"context"
	"log"
	"os"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/platform/rabbitmq"
	"github.com/danilobml/workstream/internal/workstream-notifications/mongodb"
	"github.com/danilobml/workstream/internal/workstream-notifications/readiness"
	"github.com/danilobml/workstream/internal/workstream-notifications/repositories"
	"github.com/danilobml/workstream/internal/workstream-notifications/services"
)

const (
	serviceName     = "workstream-notifications"
	httpPortName    = "NOTIFICATIONS_HTTP_PORT"
	rabbitmqUrlName = "RABBITMQ_URL"
	mongodbUriName  = "MONGODB_URI"
	mongoDbName     = "notifications"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rabbitmqUrl := os.Getenv(rabbitmqUrlName)
	if rabbitmqUrl == "" {
		log.Fatal("unable to read RABBITMQ_URL from env")
	}

	mongodbUri := os.Getenv(mongodbUriName)
	if mongodbUri == "" {
		log.Fatal("unable to read MONGODB_URI from env")
	}

	rabbitClient, err := rabbitmq.NewRabbitMQClient(ctx, rabbitmqUrl, rabbitmq.Exchange)
	if err != nil {
		log.Fatal("workstream-notifications - failed to connect to RabbitMQ", err)
	}
	defer rabbitClient.Close()

	mongoDb, mongoClient, err := mongodb.InitMongoDB(ctx, mongodbUri, mongoDbName)
	if err != nil {
		log.Fatal("workstream-notifications - failed to connect to MongoDB", err)
	}
	defer mongoClient.Disconnect(ctx)

	if err := mongodb.ApplyDbIndexes(ctx, mongoDb); err != nil {
		log.Fatal(err)
	}

	mongoRepo := repositories.NewMongoRepo(mongoDb)

	mongoService := services.NewMongoService(mongoRepo)
	rabbitService := services.NewRabbitConsumerService(rabbitClient, mongoService)

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
