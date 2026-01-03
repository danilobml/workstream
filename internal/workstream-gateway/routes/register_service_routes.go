package routes

import (
	"fmt"
	"net/http"
)

func RegisterGatewayServiceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Test gateway route")
	})
}
