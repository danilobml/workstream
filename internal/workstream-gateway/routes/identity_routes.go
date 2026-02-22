package routes

import (
	"net/http"

	"github.com/danilobml/workstream/internal/workstream-gateway/middleware"
	"github.com/danilobml/workstream/internal/workstream-gateway/handlers"
)

func RegisterIdentityRoutes(identityHandler *handlers.IdentityHandler, auth middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	// public (open)
	mux.HandleFunc("POST /register", identityHandler.Register)
	mux.HandleFunc("POST /login", identityHandler.Login)

	// admin (protected)
	mux.Handle("GET /users", auth(http.HandlerFunc(identityHandler.GetAllUsers)))
	mux.Handle("DELETE /users/{id}", auth(http.HandlerFunc(identityHandler.UnregisterUser)))

	// TODO - add routes:
	/*
		// public (protected)
		protected.HandleFunc("POST /request-password", identityHandler.RequestPasswordReset)
		protected.HandleFunc("PUT /reset-password", identityHandler.ResetPassword)
		protected.HandleFunc("POST /check-user", identityHandler.CheckUser)

		// admin (protected)
		protected.HandleFunc("GET /users/data", identityHandler.GetUserData)
		protected.HandleFunc("PUT /users/{id}", identityHandler.UpdateUser)
		protected.HandleFunc("DELETE /users/{id}/remove", identityHandler.RemoveUser)
	*/

	return mux
}
