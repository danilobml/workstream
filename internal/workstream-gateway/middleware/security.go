package middleware

import (
	"net/http"

	"github.com/unrolled/secure"
)

func Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		secureMiddleware := secure.New(secure.Options{
			FrameDeny:          true,
			ContentTypeNosniff: true,
			BrowserXssFilter:   true,
			ReferrerPolicy:     "strict-origin-when-cross-origin",
		})

		secureMiddleware.Handler(next).ServeHTTP(w, r)
	})
}
