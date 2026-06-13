CREATE TABLE audit_logs (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id          UUID            NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id             UUID            REFERENCES users(id) ON DELETE SET NULL, 

    action              VARCHAR(50)     NOT NULL, 
    entity_type         VARCHAR(50)     NOT NULL, 
    entity_id           UUID,
    changes             JSONB,
    ip_address          VARCHAR(45),
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);