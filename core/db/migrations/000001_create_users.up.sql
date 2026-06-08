CREATE TABLE users (
    wallet_address VARCHAR(42) PRIMARY KEY,
    display_name VARCHAR(100),
    bio TEXT,
    avatar_url TEXT,
    twitter_id VARCHAR(100),
    twitter_handle VARCHAR(100),
    discord_id VARCHAR(100),
    discord_handle VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);