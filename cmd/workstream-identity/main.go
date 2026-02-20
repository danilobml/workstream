package main

import (
	"context"
	"log"
	"os"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/platform/jwt"
	"github.com/danilobml/workstream/internal/platform/rabbitmq"
	"github.com/danilobml/workstream/internal/workstream-identity/db"
	identitygrpc "github.com/danilobml/workstream/internal/workstream-identity/grpc"
	"github.com/danilobml/workstream/internal/workstream-identity/readiness"
	"github.com/danilobml/workstream/internal/workstream-identity/repositories"
	"github.com/danilobml/workstream/internal/workstream-identity/services"
)

const (
	serviceName            = "workstream-identity"
	httpPortName           = "IDENTITY_HTTP_PORT"
	grpcPortName           = "IDENTITY_GRPC_PORT"
	postgresIdentityDsn    = "POSTGRES_IDENTITY_DSN"
	rabbitmqUrlName        = "RABBITMQ_URL"
	secretJwtKeyName       = "SECRET_JWT_KEY"
	baseRestorePassUrlName = "BASE_RESTORE_PASS_URL"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpPort := os.Getenv(httpPortName)
	if httpPort == "" {
		log.Fatal("unable to read IDENTITY_HTTP_PORT from env")
	}

	grpcPort := os.Getenv(grpcPortName)
	if grpcPort == "" {
		log.Fatal("unable to read IDENTITY_GRPC_PORT from env")
	}

	postgresDsn := os.Getenv(postgresIdentityDsn)
	if postgresDsn == "" {
		log.Fatal("unable to read POSTGRES_IDENTITY_DSN from env")
	}

	rabbitmqUrl := os.Getenv(rabbitmqUrlName)
	if rabbitmqUrl == "" {
		log.Fatal("unable to read RABBITMQ_URL from env")
	}

	secretJwtKey := os.Getenv(secretJwtKeyName)
	if secretJwtKey == "" {
		log.Fatal("unable to read SECRET_JWT_KEY from env")
	}

	baseRestorePassUrl := os.Getenv(baseRestorePassUrlName)
	if baseRestorePassUrl == "" {
		log.Fatal("unable to read BASE_RESTORE_PASS_URL from env")
	}

	listener, err := identitygrpc.StartGrpcListener(grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	dbConnPool, err := db.InitDB(postgresDsn)
	if err != nil {
		log.Fatal("workstream-identity - failed to initialize database", err)
	}
	defer dbConnPool.Close()

	rabbitClient, err := rabbitmq.NewRabbitMQClient(ctx, rabbitmqUrl, rabbitmq.NotificationsExchange)
	if err != nil {
		log.Fatal("workstream-identity - failed to connect to RabbitMQ", err)
	}
	defer rabbitClient.Close()

	repo := repositories.NewUserPgRepository(dbConnPool)

	jwtManager := jwt.NewJwtManager([]byte(secretJwtKey))
	rabbitService := services.NewRabbitProducerService(rabbitClient)
	userService := services.NewUserService(repo, jwtManager, rabbitService, baseRestorePassUrl)

	errCh := make(chan error, 1)
	go identitygrpc.RegisterGrpcServer(userService, listener, jwtManager, errCh)
	go func() {
		if err := <-errCh; err != nil {
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
