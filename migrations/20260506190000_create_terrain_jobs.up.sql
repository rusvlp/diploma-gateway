CREATE TABLE IF NOT EXISTS terrain_jobs (
    id          UUID PRIMARY KEY,
    user_id     UUID        NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    model_url   TEXT        NOT NULL DEFAULT '',
    texture_url TEXT        NOT NULL DEFAULT '',
    error       TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_terrain_jobs_user_id ON terrain_jobs(user_id);
