package main

import (
	"log"
	"os"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/workstream-gateway/grpc"
	"github.com/danilobml/workstream/internal/workstream-gateway/handlers"
	"github.com/danilobml/workstream/internal/workstream-gateway/readiness"
	"github.com/danilobml/workstream/internal/workstream-gateway/routes"
	services "github.com/danilobml/workstream/internal/workstream-gateway/services/adapters"
)

const (
	serviceName  = "workstream-gateway"
	httpPortName = "GATEWAY_HTTP_PORT"
	grpcAddrName = "TASKS_GRPC_ADDR"
)

func main() {
	grpcAddr := os.Getenv(grpcAddrName)

	conn, err := grpc.CreateGrpcClient(grpcAddr)
	grpc.SetClient(conn, err)
	if err != nil {
		log.Fatal(err)
	}

	service := services.NewTasksServiceClient(conn)
	gatewayHandler := handlers.NewGatewayHandler(service)

	if err := http.StartServer(
		serviceName,
		httpPortName,
		routes.RegisterGatewayServiceRoutes(gatewayHandler),
		readiness.IsReady,
	); err != nil {
		log.Fatal(err)
	}
}
