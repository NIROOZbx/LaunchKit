CREATE TABLE project_members (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID        NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    role        VARCHAR(20) NOT NULL DEFAULT 'viewer' 
                            CHECK (role IN ('owner', 'admin', 'campaign_manager', 'viewer')),

    status      VARCHAR(20) NOT NULL DEFAULT 'invited' 
                            CHECK (status IN ('invited', 'active', 'suspended')),

    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_project_user UNIQUE (project_id, user_id)
);