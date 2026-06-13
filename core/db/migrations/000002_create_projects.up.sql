CREATE TABLE projects(
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),

    company_name        VARCHAR(100)    NOT NULL,
    employee_count      VARCHAR(50),
    discovery_source    VARCHAR(100),

    name                VARCHAR(100)    NOT NULL,
    slug                VARCHAR(100)    NOT NULL UNIQUE,
    description         TEXT,
    logo_url            TEXT,
    website_url         TEXT,
    twitter_handle      VARCHAR(100),
    discord_invite_link TEXT,

    treasury_wallet     VARCHAR(42)     NOT NULL,
    blockchain          VARCHAR(50)     NOT NULL,
    environment         VARCHAR(20)     NOT NULL,
    token_address       VARCHAR(42)     NOT NULL,
    token_name          VARCHAR(100),
    token_symbol        VARCHAR(20),

    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_blockchain CHECK (blockchain IN ('ethereum', 'base', 'arbitrum')),
    CONSTRAINT chk_environment CHECK (environment IN ('mainnet', 'testnet'))
);