package routes

import (
	"net/http"

	"github.com/danilobml/workstream/internal/workstream-gateway/middleware"
	"github.com/danilobml/workstream/internal/workstream-gateway/handlers"
)

func RegisterIdentityRoutes(identityHandler *handlers.IdentityHandler, auth middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	// public
	mux.HandleFunc("POST /register", identityHandler.Register)
	mux.HandleFunc("POST /login", identityHandler.Login)
	mux.HandleFunc("POST /request-password", identityHandler.RequestPasswordReset)
	//mux.HandleFunc("PUT /reset-password", identityHandler.ResetPassword)

	// protected
	mux.Handle("PATCH /users/{id}/unregister", auth(http.HandlerFunc(identityHandler.UnregisterUser)))

	// admin
	mux.Handle("GET /users", auth(http.HandlerFunc(identityHandler.GetAllUsers)))
	mux.Handle("DELETE /users/{id}", auth(http.HandlerFunc(identityHandler.RemoveUser)))

	// TODO - add routes:
	/*
		protected.HandleFunc("GET /users/{id}", identityHandler.GetUser)
		protected.HandleFunc("PUT /users/{id}", identityHandler.UpdateUser)
		
	*/

	return mux
}
