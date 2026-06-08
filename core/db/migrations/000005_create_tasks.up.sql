CREATE TABLE tasks (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id  UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    title        VARCHAR(200) NOT NULL,
    description  TEXT,
    task_type    VARCHAR(50) NOT NULL,    -- e.g. hold_token, wallet_age, follow_twitter
    config       JSONB NOT NULL DEFAULT '{}',
    points       INT NOT NULL DEFAULT 0,
    is_required  BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order   INT NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_campaign ON tasks(campaign_id);