package httpx

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/danilobml/workstream/internal/platform/routes"
)

func StartServer(serviceName, portName string, registerServiceRoutes func(*http.ServeMux), isReady func() error) error {
	port := os.Getenv(portName)
	if port == "" {
		msg := fmt.Sprintf("%s: %s variable could not be retrieved from env", serviceName, portName)
		return errors.New(msg)
	}

	addr := fmt.Sprintf(":%s", port)

	mux := http.NewServeMux()
	routes.RegisterHealthRoutes(mux, serviceName, isReady)

	if registerServiceRoutes != nil {
		registerServiceRoutes(mux)
	}

	log.Printf("%s listening on port %s...", serviceName, port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		return err
	}

	return nil
}
