package routes

import (
	"fmt"
	"net/http"
)

func RegisterHealthRoutes(mux *http.ServeMux, serviceName string, isReady func() error) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "%s - healthz ok", serviceName)
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, req *http.Request) {
		if isReady == nil {
			fmt.Fprintf(w, "%s - readyz ok", serviceName)
			return
		}
		
		err := isReady()
		if err != nil {
			http.Error(w, "error: " + err.Error(), http.StatusServiceUnavailable)
			return
		}

		fmt.Fprintf(w, "%s - readyz ok", serviceName)
	})
}
