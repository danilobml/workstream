package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

func RateLimit(next http.Handler) http.Handler {
	limiter := httprate.NewRateLimiter(
		100,          // requests
		1*time.Minute, // per duration
		httprate.WithKeyFuncs(
			httprate.KeyByIP,
		),
	)

	return limiter.Handler(next)
}
