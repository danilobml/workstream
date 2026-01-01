# Workstream — Architecture Overview

## 1. Purpose
Workstream will be a small, production-ready microservices system that tracks work items (“tasks”) and emits events as they progress through their lifecycle.

This project is intentionally designed to be minimal, not feature-complete and optimized for learning:
- service boundaries
- synchronous (REST/gRPC) and asynchronous (event-driven) communication
- multiple persistence (Postgres + MongoDB) plus Redis
- containerized deployment (local Docker Compose, Kubernetes)

---

## 2. System Components

### 2.1 Public API Gateway
**Responsibility**
- Exposes HTTP REST API to external clients
- Performs request validation and routing
- Translates REST requests into internal gRPC calls

**Does NOT**
- Own business logic
- Persist data

**Communication**
- REST (client-service)
- gRPC client to `workstream-tasks` (service-service)

---

### 2.2 Tasks Service (Core)
**Responsibility**
- Owns domain logic and persistence for tasks
- Validates domain invariants
- Emits domain events after state changes

**Data**
- Postgres (single-service ownership)

**Communication**
- gRPC server for internal calls (from gateway)
- Event publisher to RabbitMQ

---

### 2.3 Notifications Service (Async Processor)
**Responsibility**
- Consumes task domain events
- Performs a single scoped side-effect: records processed notifications for observability/debugging
- Guarantees idempotent processing of events

**Data**
- MongoDB for processed event records and payload snapshots
- Redis for idempotency (deduplication by `event_id`)

**Communication**
- Event consumer from RabbitMQ
- REST (read-only debug API)

---

## 3. Communication Contracts

### 3.1 REST API (External)
All external traffic enters via the gateway.

| Endpoint | Purpose |
|--------|--------|
| POST /api/v1/tasks | Create task |
| GET /api/v1/tasks | List tasks |
| POST /api/v1/tasks/{id}/complete | Complete task |

---

### 3.2 gRPC API (Internal)
Gateway calls Tasks via gRPC.

| RPC | Purpose |
|----|--------|
| CreateTask | Create task |
| GetTask | Retrieve task |
| ListTasks | List tasks |
| CompleteTask | Mark task as completed |

---

### 3.3 Events (RabbitMQ)

| Event | Producer | Consumer |
|-----|---------|----------|
| workstream.task.created.v1 | workstream-tasks | workstream-notifications |
| workstream.task.completed.v1 | workstream-tasks | workstream-notifications |

**Event envelope**
- event_id
- event_type
- occurred_at
- producer
- trace_id
- payload

**RabbitMQ naming**
- Exchange: `workstream.events`
- Queue: `workstream.notifications`
- Dead-letter queue: `workstream.notifications.dlq`

---

## 4. Persistence Strategy

### 4.1 Tasks Service (Postgres)
- Relational source of truth for tasks
- Schema-first design
- Only the Tasks service accesses this database

### 4.2 Notifications Service (MongoDB)
- Append-only storage of processed event records
- Payload snapshots retained for debugging and traceability

### 4.3 Redis
- Purpose: consumer idempotency
- Key pattern: `workstream:processed:<event_id>`
- TTL: 24 hours

---

## 5. Reliability & Failure Handling
- Message delivery: at-least-once
- Consumer idempotency via Redis
- Dead-letter queue for poison messages
- Timeouts and retries configured for gRPC calls and event processing

---

## 6. Configuration & Secrets
- Configuration via environment variables
- Secrets managed outside the repository
- No reliance on implicit runtime defaults

---

## 7. Observability
- Structured logs (service, request_id, trace_id)
- Health endpoint per service
- Readiness checks include database and queue connectivity

---

## 8. Deployment Model

### 8.1 Local Development
- Docker Compose for all dependencies
- One-command startup for the entire system

### 8.2 Kubernetes (Learning Deployment)
- One Deployment per service
- One Service per Deployment
- Ingress exposes the gateway only
- No shared volumes or databases between services

---

## 9. Non-Goals
- Multi-region deployment
- Autoscaling beyond basic replicas
- End-user auth (OAuth/JWT) (can be added later)
- Exactly-once messaging (we use at-least-once + idempotency)

---

## 10. Definition of Done
- System starts locally with one command
- End-to-end flow verified (task → event → notification)
- Kubernetes manifests deploy successfully to a local cluster
- Contracts documented and versioned
- No service depends on another service’s database
