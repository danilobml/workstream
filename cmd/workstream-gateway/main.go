package main

import (
	"log"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/workstream-gateway/routes"
	"github.com/danilobml/workstream/internal/workstream-gateway/readiness"
)

const (
	serviceName = "workstream-gateway"
	portName    = "GATEWAY_HTTP_PORT"
)

func main() {
	if err := http.StartServer(
		serviceName,
		portName,
		routes.RegisterGatewayServiceRoutes,
		readiness.IsReady,
		); err != nil {
		log.Fatal(err)
	}
}
