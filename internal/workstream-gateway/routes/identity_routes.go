package routes

import (
	"net/http"

	"github.com/danilobml/workstream/internal/platform/httpx/middleware"
	"github.com/danilobml/workstream/internal/workstream-gateway/handlers"
)

func RegisterIdentityRoutes(identityHandler *handlers.IdentityHandler, auth middleware.Middleware) http.Handler {
	mux := http.NewServeMux()

	// public
	mux.HandleFunc("POST /register", identityHandler.Register)
	mux.HandleFunc("POST /login", identityHandler.Login)

	protected := http.NewServeMux()

	// TODO - add protected routes:
	/*
		protected.HandleFunc("POST /request-password", identityHandler.RequestPasswordReset)
		protected.HandleFunc("PUT /reset-password", identityHandler.ResetPassword)
		protected.HandleFunc("POST /check-user", identityHandler.CheckUser)

		protected.HandleFunc("GET /users/data", identityHandler.GetUserData)
		protected.HandleFunc("DELETE /users/{id}", identityHandler.UnregisterUser)
		protected.HandleFunc("PUT /users/{id}", identityHandler.UpdateUser)
		protected.HandleFunc("GET /users", identityHandler.GetAllUsers)
		protected.HandleFunc("DELETE /users/{id}/remove", identityHandler.RemoveUser)
	*/
	
	mux.Handle("/users/", auth(http.StripPrefix("/users", protected)))

	return mux
}
