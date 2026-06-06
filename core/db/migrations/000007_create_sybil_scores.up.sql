CREATE TABLE sybil_scores (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id           UUID NOT NULL REFERENCES campaigns(id),
    wallet_address        VARCHAR(42) NOT NULL REFERENCES users(wallet_address),
    wallet_age_score      INT NOT NULL DEFAULT 0,
    tx_count_score        INT NOT NULL DEFAULT 0,
    speed_score           INT NOT NULL DEFAULT 0,
    passport_score        INT NOT NULL DEFAULT 0,
    funding_source_score  INT NOT NULL DEFAULT 0,
    ip_cluster_score      INT NOT NULL DEFAULT 0,
    total_score           INT NOT NULL DEFAULT 0,
    flags                 TEXT[] NOT NULL DEFAULT '{}',
    is_eligible           BOOLEAN NOT NULL DEFAULT TRUE,
    computed_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(campaign_id, wallet_address)
);

CREATE INDEX idx_sybil_campaign          ON sybil_scores(campaign_id);
CREATE INDEX idx_sybil_eligible          ON sybil_scores(campaign_id, is_eligible);