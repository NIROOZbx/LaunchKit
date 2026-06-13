CREATE TYPE campaign_status AS ENUM (
    'draft', 'active', 'paused', 'ended',
    'distributing', 'completed', 'cancelled'
);

CREATE TYPE reward_type AS ENUM (
    'flat', 'points', 'tiered'
);

CREATE TYPE vesting_type AS ENUM (
    'instant', 'linear'
);

CREATE TABLE campaigns (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id          UUID            NOT NULL REFERENCES projects(id),
    title               VARCHAR(255)    NOT NULL,
    slug                VARCHAR(100)    UNIQUE NOT NULL,
    description         TEXT,
    banner_url          TEXT,
    chain               VARCHAR(20)     NOT NULL,
    token_contract      VARCHAR(42)     NOT NULL,
    token_symbol        VARCHAR(20),
    token_decimals      INT             DEFAULT 18,
    total_allocation    NUMERIC(78, 0)  NOT NULL,
    claimed_amount      NUMERIC(78, 0)  NOT NULL DEFAULT 0,
    reward_type         reward_type     NOT NULL DEFAULT 'points',
    reward_config       JSONB           NOT NULL DEFAULT '{}',
    eligibility_rules   JSONB           NOT NULL DEFAULT '{}',
    start_date          TIMESTAMPTZ     NOT NULL,
    end_date            TIMESTAMPTZ     NOT NULL,
    claim_window_days   INT             NOT NULL DEFAULT 30,
    vesting_type        vesting_type    NOT NULL DEFAULT 'instant',
    vesting_days        INT             NOT NULL DEFAULT 0,
    gas_sponsored       BOOLEAN         NOT NULL DEFAULT FALSE,
    merkle_root         VARCHAR(66),
    claim_contract      VARCHAR(42),
    deploy_tx_hash      VARCHAR(66),
    status              campaign_status NOT NULL DEFAULT 'draft',
    finalized_at        TIMESTAMPTZ,
    completed_at        TIMESTAMPTZ,
    cancelled_at        TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ     DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     DEFAULT NOW()
);

-- constraints
ALTER TABLE campaigns ADD CONSTRAINT campaigns_end_after_start
    CHECK (end_date > start_date);
ALTER TABLE campaigns ADD CONSTRAINT campaigns_min_claim_window
    CHECK (claim_window_days >= 1);
ALTER TABLE campaigns ADD CONSTRAINT campaigns_vesting_consistency
    CHECK (
        (vesting_type = 'instant' AND vesting_days = 0) OR
        (vesting_type = 'linear'  AND vesting_days > 0)
    );
ALTER TABLE campaigns ADD CONSTRAINT campaigns_positive_allocation
    CHECK (total_allocation > 0);
ALTER TABLE campaigns ADD CONSTRAINT campaigns_claimed_within_total
    CHECK (claimed_amount <= total_allocation);

-- indexes
CREATE INDEX idx_campaigns_project_id   ON campaigns(project_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_campaigns_status       ON campaigns(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_campaigns_chain        ON campaigns(chain) WHERE deleted_at IS NULL;
CREATE INDEX idx_campaigns_slug         ON campaigns(slug);
CREATE INDEX idx_campaigns_dates        ON campaigns(start_date, end_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_campaigns_active_chain ON campaigns(chain, start_date, end_date) WHERE status = 'active' AND deleted_at IS NULL;

-- trigger
CREATE TRIGGER campaigns_updated_at
    BEFORE UPDATE ON campaigns
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();