CREATE TABLE auth_nonces (
    wallet_address  VARCHAR(42)     PRIMARY KEY,
    nonce           VARCHAR(64)     NOT NULL,
    expires_at      TIMESTAMPTZ     NOT NULL,
    created_at      TIMESTAMPTZ     DEFAULT NOW()
);

CREATE INDEX idx_auth_nonces_expires
    ON auth_nonces(expires_at);