package routes

import (
	"net/http"

	"github.com/danilobml/workstream/api/openapi"
)

func RegisterOpenapiRoutes(root *http.ServeMux) {
	// Yaml
	root.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		_, _ = w.Write(openapi.SpecYAML())
	})

	// Swagger UI
	root.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(openapi.SwaggerHTML())
	})
	root.HandleFunc("GET /docs/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(openapi.SwaggerHTML())
	})
}
