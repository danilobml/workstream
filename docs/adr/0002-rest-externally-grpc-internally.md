# ADR 0002: REST externally, gRPC internally

Status: Accepted

## Decision
The setup exposes a **REST/HTTP API** to external clients via `gateway`, and uses **gRPC** for internal service-to-service calls (e.g., gateway → `tasks`).

## Context
External clients need a simple, widely compatible interface that. Internally, services benefit from a stricter contract and efficient communication.

## Consequences
- Only the gateway is exposed to external clients (via HTTP).
- Internal calls use gRPC and are not publicly exposed.
- REST and gRPC contracts must stay consistent (also naming).
