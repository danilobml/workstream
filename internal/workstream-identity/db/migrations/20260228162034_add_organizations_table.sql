-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS organization_id UUID;

-- Creates a default org and assigns it to existing users:
WITH org AS (
    INSERT INTO organizations (name)
    VALUES ('default')
    ON CONFLICT DO NOTHING
    RETURNING id
)
UPDATE users
SET organization_id = COALESCE(
    users.organization_id,
    (SELECT id FROM org),
    (SELECT id FROM organizations WHERE name = 'default' LIMIT 1)
)
WHERE users.organization_id IS NULL;

-- NOT NULL and FK
ALTER TABLE users
    ALTER COLUMN organization_id SET NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT fk_users_organization_id
    FOREIGN KEY (organization_id)
    REFERENCES organizations(id)
    ON DELETE RESTRICT;

-- Index
CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_organization_id;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS fk_users_organization_id;

ALTER TABLE users
    DROP COLUMN IF EXISTS organization_id;

DROP TABLE IF EXISTS organizations;
-- +goose StatementEnd