CREATE TYPE completion_status AS ENUM (
    'pending',  -- Needs manual review by admin
    'verified', -- Success, points awarded
    'rejected'  -- Failed or denied by admin
);

CREATE TABLE tasks (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id         UUID            NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,

    title               VARCHAR(255)    NOT NULL,
    description         TEXT,

    -- 1. FRONTEND TAG: Tells the UI what icon/component to render.
    -- (e.g., 'tiktok_follow', 'upload_meme', 'custom_game_check')
    task_type           VARCHAR(100)    NOT NULL,

    -- 2. BACKEND ROUTER: Tells the Go backend how to verify it.
    -- (e.g., 'manual', 'webhook_api', 'twitter_api', 'web3_rpc')
    verification_type   VARCHAR(100)    NOT NULL,

    -- 3. THE PAYLOAD: The exact rules for the verification_type.
    config              JSONB           NOT NULL DEFAULT '{}',

    points              INT             NOT NULL DEFAULT 100,
    is_required         BOOLEAN         NOT NULL DEFAULT FALSE,
    display_order       INT             NOT NULL DEFAULT 0,

    deleted_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ     DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     DEFAULT NOW()
);