package routes

import (
	"net/http"

	"github.com/danilobml/workstream/internal/workstream-gateway/handlers"
	"github.com/danilobml/workstream/internal/workstream-gateway/httpx/middleware"
)

func RegisterGatewayServiceRoutes(handler *handlers.GatewayHandler) func(mux *http.ServeMux) {

	return func(mux *http.ServeMux) {
		mux.HandleFunc("POST /tasks", handler.CreateNewTask)
		mux.HandleFunc("GET /tasks", handler.GetTasks)
		mux.HandleFunc("GET /tasks/{id}", handler.GetTask)
		mux.HandleFunc("POST /tasks/{id}/complete", handler.CompleteTask)

		// Global middleware
		use := middleware.ApplyMiddlewares(
			middleware.Recover,
			middleware.RequestId,
			middleware.Logger,
			middleware.Dos,
			middleware.RateLimit,
			middleware.Cors,
			middleware.Security,
		)

		use(mux)
	}
}
