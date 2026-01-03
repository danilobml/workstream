package routes

import (
	"fmt"
	"net/http"
)

func RegisterHealthRoutes(mux *http.ServeMux, serviceName string) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s - healthz ok", serviceName)
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s - readyz ok", serviceName)
	})
}
