package main

import (
	"log"

	http "github.com/danilobml/workstream/internal/platform/httpserver"
)

const (
	serviceName = "workstream-tasks"
	portName    = "TASKS_HTTP_PORT"
)

func main() {
	if err := http.StartServer(serviceName, portName, nil); err != nil {
		log.Fatal(err)
	}
}
