# ADR-0008: Platform Admin (User Property) and Organization Routes

- **Status:** Accepted  
- **Date:** 2026-03-22  

## Context

Workstream is evolving into a multi-tenant system where data must be isolated per organization (tenant).

We need:

- Organization endpoints to manage organization lifecycle.
- A platform-level administrator capable of managing organizations across tenants.
- A clear separation between organization-scoped roles and global privileges.

## Decision

### 1. Platform Admin as User Property

Add a boolean property `is_platform_admin` to the `users` table.

Platform admins are global "superusers" who:

- Can create organizations.
- Can operate across all organizations.
- Are not restricted to a single tenant.

This property is **not modeled as a role**, because roles are scoped to organizations, while platform admin is a **cross-tenant capability**.

### 2. Platform Admin Only Added Directly in the DB

Platform admins are bootstrap users and can only be created:

- Directly via database migration
- Or manual DB insertion

No API route will exist to grant or revoke platform admin privileges.

This prevents privilege escalation via API misuse.

### 3. Organization Routes Restricted to Platform Admins

Organization management endpoints are accessible only to platform admins.

Identity issues JWTs containing at minimum:

- `sub` (user_id)
- `org_id` (organization_id)
- `roles` (organization-scoped roles)
- `is_platform_admin` (boolean)

Authorization logic:

- If `is_platform_admin == true`, allow cross-tenant operations
- Otherwise enforce `org_id` scoping

## Consequences

- Introduces global authorization capability
- Simplifies cross-tenant operations
- Requires JWT claim extension
- Requires new authorization middleware
- Adds boolean property to `users` table
- Platform admin bootstrap must be handled manually
- Future audit logging should track platform-admin actions

## Alternatives Considered

### 1. Use a Role to Define Platform Admin

Rejected because:

- Roles are organization-scoped
- Platform admin must be global
- Avoids mixing tenant and platform authorization concerns
- Simplifies authorization checks

## Follow-ups

- Implement migrations:
  - Identity: `organizations` table
  - Identity: `users.is_platform_admin BOOLEAN DEFAULT false`
- Update Identity JWT issuance to include `is_platform_admin`
- Create Organization routes
- Add Gateway authorization middleware for platform-admin-only operations
- Add integration tests for cross-tenant authorization
