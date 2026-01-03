package main

import (
	"log"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/workstream-notifications/readiness"
)

const (
	serviceName = "workstream-notifications"
	portName    = "NOTIFICATIONS_HTTP_PORT"
)

func main() {
	if err := http.StartServer(
		serviceName,
		portName,
		nil,
		readiness.IsReady,
		); err != nil {
		log.Fatal(err)
	}
}
