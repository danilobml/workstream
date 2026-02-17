package routes

import (
	"net/http"

	"github.com/danilobml/workstream/internal/workstream-gateway/handlers"
)

func RegisterTaskRoutes(handler *handlers.TasksHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", handler.CreateNewTask)
	mux.HandleFunc("GET /", handler.GetTasks)
	mux.HandleFunc("GET /{id}", handler.GetTask)
	mux.HandleFunc("POST /{id}/complete", handler.CompleteTask)

	return mux
}
