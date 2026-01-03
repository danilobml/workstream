# ADR 0003: Use postgres for tasks service and mongoDB for notifications, plus Redis

Status: Accepted

## Decision
Tasks stores tasks in Postgres; Notifications stores processed event records in MongoDB; Redis is used only for event deduplication by event_id.

What each service uses:
- Tasks: **Postgres**
- Notifications: **MongoDB**
- **Redis** for idempotency/dedupe only

Rule: no shared databases across services.

## Context
- Tasks data is relational + needs constraints/transactions -> Postgres fits.
- Notifications stores processed event records / snapshots (append-only, flexible schema) -> Mongo fits.
- Redis is ideal for fast TTL-based dedupe keys.

## Consequences
- Service-owned data: only Tasks touches Postgres; only Notifications touches Mongo.
- No cross-service reads. Must use gRPC/events.
- Two backup/restore stories (Postgres + Mongo) even in local dev.
- Migration tooling differs: Postgres uses schema migrations; Mongo uses collection/index setup
- Redis keys must have a TTL and be keyed by event_id

## Alternatives
- Single database for everything (simpler, but violates boundaries/learning goals)
- Postgres for both services (simpler ops, less flexible event record storage)
- Kafka instead of RabbitMQ (heavier for this project; not needed)
