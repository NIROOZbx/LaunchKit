CREATE TABLE users (
    -- Web3 Core Identity (Primary Key)
    wallet_address      VARCHAR(42) PRIMARY KEY, 
    ens_name            VARCHAR(255),

    display_name        VARCHAR(100),
    avatar_url          TEXT,

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

-- Index 1: The B2C Social Verification Acceleration Loops
CREATE INDEX idx_users_twitter_id ON users(twitter_id) WHERE twitter_id IS NOT NULL;
CREATE INDEX idx_users_discord_id ON users(discord_id) WHERE discord_id IS NOT NULL;

-- Index 2: System Session & Security Audit Index
CREATE INDEX idx_users_last_seen ON users(last_seen DESC);

-- Index 3: Case-Insensitive B2B Team Workspace Lookup Bridge
CREATE INDEX idx_users_ens_lower ON users(LOWER(ens_name)) WHERE ens_name IS NOT NULL;