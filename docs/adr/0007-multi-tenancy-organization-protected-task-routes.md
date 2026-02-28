# ADR-0007: Organization (Tenant) Model and JWT-Protected Task Routes

- **Status:** Accepted
- **Date:** 2026-02-20

## Context
Workstream is evolving into a multi-tenant system where data must be isolated per organization (tenant). We need:
- A first-class **Organization** entity in Identity.
- Users to belong to an organization at creation time.
- Tasks to be scoped to an organization.
- All task endpoints to require authentication, and to enforce tenant isolation without trusting client-provided organization identifiers.

## Decision
1. **Identity owns Organizations**
   - Add an `organizations` table to the existing `postgres-identity` database.
   - Add `users.organization_id` (FK to `organizations.id`), set at user creation (registration).

2. **JWT includes tenant context**
   - Identity issues JWTs containing at minimum: `sub` (user_id) and `org_id` (organization_id), plus `roles`.
   - Services do not accept `org_id` from request body/query for tenant scoping.

3. **Tasks are organization-scoped**
   - Add `tasks.organization_id` to the Tasks service database schema.
   - All task reads/writes must filter by `organization_id` derived from JWT claims (via gateway context).

4. **Protect all task routes**
   - All `/tasks` routes in the Gateway require valid JWT authentication.
   - Gateway extracts `org_id` from JWT and passes it to Tasks service calls.
   - Tasks service methods require `org_id` and enforce `WHERE organization_id = $org_id` for all queries.

## Consequences
- Strong tenant isolation by default (no cross-org access via IDs alone).
- Clear ownership: Identity manages org membership; Tasks stores org-scoped data.
- Simplified API usage: clients never send org identifiers for task operations.
- Requires schema migrations in Identity and Tasks, plus updates to JWT claims and gateway middleware.
- Backfill needed for any existing users/tasks (assign to a default or migrated organization).

## Alternatives Considered
1. **Task service calls Identity on each request to resolve org**
   - Rejected: adds latency, coupling, and failure modes.

2. **Client provides org_id in requests**
   - Rejected: insecure; easy to bypass tenant isolation.

3. **Single global DB with row-level security (RLS)**
   - Deferred: higher complexity; can be considered later for defense-in-depth.

## Follow-ups
- Implement migrations:
  - Identity: `organizations` + `users.organization_id`
  - Tasks: `tasks.organization_id` + indexes
- Update Identity JWT issuance to include `org_id`.
- Add Gateway JWT auth middleware for all task routes.
- Update Tasks service APIs and queries to require `org_id` for all operations.
