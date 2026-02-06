# ADR-0005: HTTP Middleware Stack for Go REST APIs

## Status
Accepted

## Context
The Workstream services expose HTTP REST APIs using Go's standard `net/http` package and `http.ServeMux`.

To address common cross-cutting concerns (security, stability, observability, and abuse protection), a consistent middleware stack is required. The middleware must be:
- Compatible with plain `net/http`
- Explicit and easy to reason about
- Lightweight and dependency-minimal
- Reusable across services

## Decision
We adopt a small, ordered HTTP middleware stack applied by wrapping the root `http.Handler` (the `ServeMux`) once at server startup.

Each middleware follows the standard Go signature:

```
func(next http.Handler) http.Handler
```

### Middleware Responsibilities

The stack consists of the following middlewares:

1. **Recover**
   - Recovers from panics
   - Prevents process crashes
   - Returns a safe HTTP 500 response

2. **RequestId**
   - Generates or propagates a request ID
   - Stores it in `context.Context`
   - Makes logs traceable across services

3. **Logger**
   - Logs request lifecycle (method, path, status, duration)
   - Includes request ID when available

4. **DoS Protection**
   - Limits request body size
   - Applies basic timeout / resource guards

5. **Rate Limiting**
   - Limits request rate per client (IP-based)
   - Implemented using `github.com/go-chi/httprate`
   - In-memory and per-container

6. **CORS**
   - Allows controlled cross-origin access
   - Handles browser preflight (`OPTIONS`) requests

7. **Security Headers**
   - Adds standard HTTP security headers
   - Implemented using `github.com/unrolled/secure`

### Middleware Order

The middleware is applied in the following order:

```
Recover
→ RequestId
→ Logger
→ DoS Protection
→ Rate Limit
→ CORS
→ Security Headers
→ Handlers
```

This order ensures:
- Panics are always recovered
- All logs include a request ID
- Cheap rejection happens early
- Security headers are applied to all responses

## Consequences

### Positive
- Clear and explicit request lifecycle
- No framework lock-in (pure `net/http`)
- Easy to test each middleware independently
- Consistent behavior across services

### Negative
- Rate limiting is per replica (not global)
- No automatic route-level middleware without additional routing logic

## Notes
- Middleware is applied once by wrapping the `ServeMux` and passing it to `http.Server{Handler: wrapped}`
- Health endpoints may optionally bypass some middleware (e.g. rate limiting)

## References
- Go `net/http` package
- `github.com/go-chi/httprate`
- `github.com/unrolled/secure`

