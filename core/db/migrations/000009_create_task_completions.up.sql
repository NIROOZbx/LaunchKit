CREATE TABLE task_completions (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID            NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id             UUID            NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    campaign_id         UUID            NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,

    status              completion_status NOT NULL DEFAULT 'pending',
    points_earned       INT             NOT NULL DEFAULT 0,

    -- Stores the user's submission (e.g., an image URL, a tx hash, or API response ID)
    proof               JSONB           NOT NULL DEFAULT '{}',

    completed_at        TIMESTAMPTZ     DEFAULT NOW(),

    UNIQUE(user_id, task_id)
);

-- Indexes
CREATE INDEX idx_tasks_campaign ON tasks(campaign_id);
CREATE INDEX idx_completions_leaderboard ON task_completions(campaign_id, user_id);