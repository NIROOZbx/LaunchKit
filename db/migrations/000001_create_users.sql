-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

CREATE TABLE users (
    --Internal System Identity
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Web3 Layer Authentication
    wallet_address      VARCHAR(42) NOT NULL UNIQUE, 
    ens_name            VARCHAR(255),

    display_name        VARCHAR(100),
    avatar_url          TEXT,

    user_type VARCHAR(20) NOT NULL DEFAULT 'b2c' CHECK (user_type IN ('b2c', 'b2b', 'admin')),

    -- Personal B2C Social Identities (For Verification Engines)
    twitter_id          VARCHAR(50)  UNIQUE,
    twitter_handle      VARCHAR(100),

    discord_id          VARCHAR(50)  UNIQUE,
    discord_handle      VARCHAR(100),
    discord_token       TEXT,
    -- Core System Infrastructure Tracking
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    last_seen           TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_wallet_address ON users(wallet_address);
-- Index 1: The B2C Social Verification Acceleration Loops
CREATE INDEX idx_users_twitter_id ON users(twitter_id) WHERE twitter_id IS NOT NULL;
CREATE INDEX idx_users_discord_id ON users(discord_id) WHERE discord_id IS NOT NULL;

-- Index 2: System Session & Security Audit Index
CREATE INDEX idx_users_last_seen ON users(last_seen DESC);

-- Index 3: Case-Insensitive B2B Team Workspace Lookup Bridge
CREATE INDEX idx_users_ens_lower ON users(LOWER(ens_name)) WHERE ens_name IS NOT NULL;

CREATE TRIGGER users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS users_updated_at ON users;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS update_updated_at();
