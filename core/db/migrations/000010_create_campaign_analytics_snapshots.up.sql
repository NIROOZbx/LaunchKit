CREATE TABLE campaign_analytics_snapshots (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id         UUID            NOT NULL REFERENCES campaigns(id),

    -- snapshot data
    participant_count   INT             DEFAULT 0,
    task_completion_rate DECIMAL(5,2)   DEFAULT 0,
    sybil_flagged_count INT             DEFAULT 0,
    claim_rate          DECIMAL(5,2)    DEFAULT 0,
    tokens_claimed      NUMERIC(78, 0)  DEFAULT 0,

    -- when
    snapshot_at         TIMESTAMPTZ     DEFAULT NOW()
);