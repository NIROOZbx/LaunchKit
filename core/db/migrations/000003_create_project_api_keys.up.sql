-- Secure Storage for Webhook API Keys
CREATE TABLE  project_api_keys (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id     UUID            NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name                VARCHAR(50)     NOT NULL DEFAULT 'Default Key',
    key_hash            VARCHAR(64)     NOT NULL UNIQUE,
    key_hint            VARCHAR(8)      NOT NULL,
    is_active           BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    last_used_at        TIMESTAMPTZ
);

CREATE INDEX idx_api_keys_hash ON organization_api_keys(key_hash);