# ADR 0001: At-least-once delivery with idempotent consumers

Status: Accepted

## Decision
Workstream uses **at-least-once** event delivery. Event consumers must be **idempotent**.

## Context
RabbitMQ and distributed systems can deliver the same message more than once due to retries or failures.

## Consequences
- Events may be processed more than once.
- Consumers deduplicate events using `event_id`.
- Redis is used to store processed `event_id`s.
- Failed messages are sent to a dead-letter queue.
