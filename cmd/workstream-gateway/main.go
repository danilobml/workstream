package main

import (
	"log"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	gatewayroutes "github.com/danilobml/workstream/internal/workstream-gateway/routes"
)

const (
	serviceName = "workstream-gateway"
	portName    = "GATEWAY_HTTP_PORT"
)

func main() {
	registerRoutesFn := gatewayroutes.RegisterGatewayServiceRoutes

	if err := http.StartServer(serviceName, portName, registerRoutesFn); err != nil {
		log.Fatal(err)
	}
}
