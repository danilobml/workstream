# ADR-0006: HTTP Middleware Stack for Go REST APIs

## Status
Accepted

## Context
The Workstream Gateway service exposes HTTP REST APIs using Go's standard `net/http` package.

To address common cross-cutting concerns (security, stability, observability, and abuse protection), a consistent middleware stack is required. The middleware must be:
- Compatible with plain `net/http`
- Reusable across services

## Decision
We adopt a small, ordered HTTP middleware stack applied by wrapping the root handler once at server startup.

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

## Consequences

### Positive
- Clear and explicit request lifecycle
- Easy to test each middleware independently
- Consistent behavior

### Negative
- Rate limiting is per replica (not global)
