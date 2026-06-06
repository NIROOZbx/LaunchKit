CREATE TYPE completion_status AS ENUM ('pending', 'verified', 'failed');

CREATE TABLE task_completions (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id    UUID NOT NULL REFERENCES campaigns(id),
    task_id        UUID NOT NULL REFERENCES tasks(id),
    wallet_address VARCHAR(42) NOT NULL REFERENCES users(wallet_address),
    status         completion_status NOT NULL DEFAULT 'pending',
    proof          JSONB NOT NULL DEFAULT '{}',
    failure_reason TEXT,
    verified_at    TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(task_id, wallet_address)  -- one completion per wallet per task
);

CREATE INDEX idx_completions_campaign_wallet ON task_completions(campaign_id, wallet_address);
CREATE INDEX idx_completions_task            ON task_completions(task_id);