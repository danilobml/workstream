package main

import (
	"log"
	"os"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/workstream-tasks/grpc"
	"github.com/danilobml/workstream/internal/workstream-tasks/readiness"
)

const (
	serviceName  = "workstream-tasks"
	httpPortName = "TASKS_HTTP_PORT"
	grpcPortName = "TASKS_GRPC_PORT"
)

func main() {
	grpcPort := os.Getenv(grpcPortName)
	if grpcPort == "" {
		log.Fatal("unable to read TASKS_GRPC_PORT from env")
	}

	listener, err := grpc.StartGrpcListener(grpcPort)
	if err != nil {
		log.Fatal(err)
	}

	errCh := make(chan error, 1)
	go grpc.RegisterGrpcServer(listener, errCh)	
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
