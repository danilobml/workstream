package main

import (
	"log"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
	"github.com/danilobml/workstream/internal/workstream-tasks/readiness"
)

const (
	serviceName = "workstream-tasks"
	portName    = "TASKS_HTTP_PORT"
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
