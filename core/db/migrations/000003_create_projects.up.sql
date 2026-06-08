CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_address VARCHAR(42) NOT NULL REFERENCES users(wallet_address),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    logo_url TEXT,
    website_url TEXT,
    twitter_url TEXT,
    discord_url TEXT,
    token_contract VARCHAR(42) NOT NULL,
    chain VARCHAR(50) NOT NULL,
    treasury_wallet VARCHAR(42) NOT NULL,
    api_key VARCHAR(100) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_owner ON projects(owner_address);