package routes

import (
	"net/http"

	"github.com/danilobml/workstream/internal/workstream-gateway/handlers"
	"github.com/danilobml/workstream/internal/workstream-gateway/middleware"
)

func RegisterIdentityRoutes(identityHandler *handlers.IdentityHandler, auth middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	// public (open)
	mux.HandleFunc("POST /register", identityHandler.Register)
	mux.HandleFunc("POST /login", identityHandler.Login)
	mux.HandleFunc("POST /request-password", identityHandler.RequestPasswordReset)
	mux.HandleFunc("PUT /reset-password", identityHandler.ResetPassword)
	// public (protected)
	mux.Handle("PATCH /users/{id}/unregister", auth(http.HandlerFunc(identityHandler.UnregisterUser)))

	// admin
	mux.Handle("GET /users", auth(http.HandlerFunc(identityHandler.GetAllUsers)))
	mux.Handle("DELETE /users/{id}", auth(http.HandlerFunc(identityHandler.RemoveUser)))
	// mux.Handle("GET /users/{id}", auth(http.HandlerFunc(identityHandler.GetUser)))
	// mux.Handle("PUT /users/{id}", auth(http.HandlerFunc(identityHandler.UpdateUser)))

	return mux
}
