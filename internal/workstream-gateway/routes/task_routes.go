package routes

import (
	"net/http"

	"github.com/danilobml/workstream/internal/workstream-gateway/handlers"
	"github.com/danilobml/workstream/internal/workstream-gateway/middleware"
)

func RegisterTaskRoutes(taskHandler *handlers.TasksHandler, auth middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	// public(protected)
	mux.Handle("POST /tasks", auth(http.HandlerFunc(taskHandler.CreateNewTask)))
	mux.Handle("GET /tasks", auth(http.HandlerFunc(taskHandler.GetTasks)))
	mux.Handle("GET /tasks/{id}", auth(http.HandlerFunc(taskHandler.GetTask)))
	mux.Handle("POST /tasks/{id}/complete", auth(http.HandlerFunc(taskHandler.CompleteTask)))

	return mux
}
