# ADR-004: User Identity Service

## Status
Accepted

## Context

Workstream needs a **user identity service** responsible for:
- User registration
- Login / authentication
- Password reset flows
- User profile data
- Token issuance and validation

The platform runs on Kubernetes and already exposes a **Workstream Gateway** as the single public HTTP entry point.

This service is designed as a **brand-new system**, with no legacy or migration constraints.

---

## Decision

### Service name

**`workstream-identity`**

Reason:
- Clearly communicates responsibility
- Broad enough to cover auth + user data
- Stable naming for future extensions (roles, orgs, policies)

---

### External API (client-facing)

**HTTP/REST via workstream-gateway**

Used for:
- `POST /register`
- `POST /login`
- `POST /request-password-reset`
- `PUT /reset-password`
- `GET /me`

Reason:
- Client-facing semantics map naturally to REST
- Gateway already handles HTTP concerns (auth, CORS, rate limiting)
- Performance difference vs gRPC is irrelevant for these flows

---

### Internal API (service-to-service)

**gRPC**

Used for:
- GetUserById
- ValidateToken
- Resolve roles / permissions
- Admin-level user queries

Reason:
- Strong typing and explicit contracts
- Low-latency internal calls
- Clear separation between external and internal APIs

---

### Database

**PostgreSQL**

Reason:
- Strong consistency guarantees
- Easy enforcement of constraints (unique email, usernames)
- Mature ecosystem and operational stability
- Ideal for authentication and identity data

---

### Token strategy

- JWT access tokens
- Short-lived tokens
- Verified locally by services when possible
- Optional gRPC validation for privileged or sensitive operations

---

## Consequences

- Workstream Gateway remains the only public HTTP surface
- Identity logic is centralized in a single service
- Internal communication is efficient and well-defined
- Service can be extended later (organizations, roles, policies)

---

## Summary

- **Service name:** workstream-identity
- **Client API:** REST (via gateway)
- **Internal API:** gRPC
- **Database:** PostgreSQL
- **Authentication:** JWT

