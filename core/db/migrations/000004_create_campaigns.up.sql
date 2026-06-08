CREATE TYPE campaign_status AS ENUM (
    'draft',
    'active',
    'paused',
    'ended',
    'distributing',
    'completed'
);

CREATE TYPE reward_type AS ENUM (
    'flat',
    'points',
    'tiered'
);

CREATE TYPE vesting_type AS ENUM (
    'instant',
    'linear'
);

CREATE TABLE campaigns (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id          UUID NOT NULL REFERENCES projects(id),
    name                VARCHAR(200) NOT NULL,
    description         TEXT,
    banner_url          TEXT,
    status              campaign_status NOT NULL DEFAULT 'draft',
    chain               VARCHAR(50) NOT NULL,
    token_contract      VARCHAR(42) NOT NULL,
    total_allocation    NUMERIC(78, 0) NOT NULL,   -- holds large uint256 values
    reward_type         reward_type NOT NULL,
    reward_config       JSONB NOT NULL DEFAULT '{}',
    eligibility_rules   JSONB NOT NULL DEFAULT '{}',
    vesting_type        vesting_type NOT NULL DEFAULT 'instant',
    vesting_days        INT,
    claim_window_days   INT NOT NULL DEFAULT 30,
    gas_sponsored       BOOLEAN NOT NULL DEFAULT FALSE,
    starts_at           TIMESTAMPTZ NOT NULL,
    ends_at             TIMESTAMPTZ NOT NULL,
    merkle_root         VARCHAR(66),               -- 0x + 64 hex chars
    contract_address    VARCHAR(42),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_campaigns_project    ON campaigns(project_id);
CREATE INDEX idx_campaigns_status     ON campaigns(status);
CREATE INDEX idx_campaigns_chain      ON campaigns(chain);