-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS organization_id UUID;

UPDATE tasks
SET organization_id = 'f0d8eb9f-952f-41fc-b803-1e7407e37a9e'
WHERE organization_id IS NULL;

ALTER TABLE tasks
    ALTER COLUMN organization_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_organization_id ON tasks(organization_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tasks_organization_id;

ALTER TABLE tasks
    DROP COLUMN IF EXISTS organization_id;
-- +goose StatementEnd
