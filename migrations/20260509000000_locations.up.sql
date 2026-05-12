CREATE TABLE IF NOT EXISTS locations (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL,
    job_id           UUID NOT NULL,
    name             VARCHAR(255) NOT NULL,
    scale_z          FLOAT NOT NULL DEFAULT 0.3,
    water_enabled    BOOLEAN NOT NULL DEFAULT FALSE,
    water_level      FLOAT NOT NULL DEFAULT 0.1,
    tree_count_percent INT NOT NULL DEFAULT 100,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_locations_user_id ON locations(user_id);
