# Workstream — Event Contracts

## Overview
This document defines the asynchronous event contracts used in Workstream.
Events are published by domain services and consumed by downstream processors.

---

## Messaging Topology

### Exchange
- Name: `workstream.events`
- Type: topic
- Purpose: fan-out domain events to interested consumers

### Queues

#### Notifications Queue
- Name: `workstream.notifications`
- Bound to exchange: `workstream.events`
- Routing keys:
  - `workstream.task.*`

#### Dead-Letter Queue
- Name: `workstream.notifications.dlq`
- Purpose: stores messages that could not be processed after retries

---

## Events

### workstream.task.created.v1
**Producer**
- workstream-tasks

**Consumers**
- workstream-notifications

**Description**
- Emitted when a task is created successfully.

- Payload contains: \task_id` (UUID).`

---

### workstream.task.completed.v1
**Producer**
- workstream-tasks

**Consumers**
- workstream-notifications

**Description**
- Emitted when a task is marked as completed.

- Payload contains: \task_id` (UUID).`

---

## Event Envelope (Common Fields)
All events share a common envelope.

- event_id
- event_type
- occurred_at
- producer
- trace_id
- payload

---

## Delivery Semantics
- Delivery guarantee: at-least-once
- Consumers must be idempotent
- Redis is used to deduplicate events by `event_id`

---

## Versioning
- Events are immutable once published
- Backward-incompatible changes require a new version suffix (v2, v3, …)
