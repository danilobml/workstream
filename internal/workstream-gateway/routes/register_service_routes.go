package routes

import (
	"net/http"

	"github.com/danilobml/workstream/internal/workstream-gateway/handlers"
)

func RegisterGatewayServiceRoutes(handler *handlers.GatewayHandler) func(mux *http.ServeMux) {
	return func(mux *http.ServeMux) {
		mux.HandleFunc("POST /tasks", handler.CreateNewTask)
		mux.HandleFunc("GET /tasks", handler.GetTasks)
		mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
		mux.HandleFunc("POST /tasks/{id}/complete", handler.CompleteTask)
	}
}
