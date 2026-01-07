package main

import (
	"log"
	"os"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/workstream-tasks/db"
	"github.com/danilobml/workstream/internal/workstream-tasks/grpc"
	"github.com/danilobml/workstream/internal/workstream-tasks/readiness"
	"github.com/danilobml/workstream/internal/workstream-tasks/repositories"
	"github.com/danilobml/workstream/internal/workstream-tasks/services"
)

const (
	serviceName  = "workstream-tasks"
	httpPortName = "TASKS_HTTP_PORT"
	grpcPortName = "TASKS_GRPC_PORT"
	postgresDsnName = "POSTGRES_DSN"
)

func main() {
	grpcPort := os.Getenv(grpcPortName)
	if grpcPort == "" {
		log.Fatal("unable to read TASKS_GRPC_PORT from env")
	}

	postgresDsn := os.Getenv(postgresDsnName)
	if postgresDsn == "" {
		log.Fatal("unable to read POSTGRES_DSN from env")
	}

	listener, err := grpc.StartGrpcListener(grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	dbConnPool, err := db.InitDB(postgresDsn)
	if err != nil {
		log.Fatal("workstream-tasks - failed to initialize database")
	}
	defer dbConnPool.Close()

	repo := repositories.NewPgTaskRepository(dbConnPool)
	tasksServer := services.NewTasksService(repo)

	errCh := make(chan error, 1)
	go grpc.RegisterGrpcServer(tasksServer, listener, errCh)	
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
