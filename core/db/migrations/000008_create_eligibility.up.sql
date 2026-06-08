CREATE TYPE claim_status AS ENUM ('unclaimed', 'pending', 'claimed', 'failed');

CREATE TABLE eligibility (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id      UUID NOT NULL REFERENCES campaigns(id),
    wallet_address   VARCHAR(42) NOT NULL REFERENCES users(wallet_address),
    is_eligible      BOOLEAN NOT NULL DEFAULT FALSE,
    total_points     INT NOT NULL DEFAULT 0,
    token_amount     NUMERIC(78, 0) NOT NULL DEFAULT 0,
    merkle_proof     JSONB NOT NULL DEFAULT '[]',
    claim_status     claim_status NOT NULL DEFAULT 'unclaimed',
    claim_tx_hash    VARCHAR(66),
    claimed_at       TIMESTAMPTZ,
    computed_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(campaign_id, wallet_address)
);

CREATE INDEX idx_eligibility_campaign        ON eligibility(campaign_id);
CREATE INDEX idx_eligibility_claim_status    ON eligibility(campaign_id, claim_status);