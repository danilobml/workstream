package main

import (
	"log"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
)

const (
	serviceName = "workstream-notifications"
	portName    = "NOTIFICATIONS_HTTP_PORT"
)

func main() {
	if err := http.StartServer(serviceName, portName, nil); err != nil {
		log.Fatal(err)
	}
}
