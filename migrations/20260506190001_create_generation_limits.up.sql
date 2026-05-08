CREATE TABLE IF NOT EXISTS generation_limits (
    id          UUID PRIMARY KEY,
    user_id     UUID    NOT NULL UNIQUE,
    daily_limit INT     NOT NULL DEFAULT 10,
    used_today  INT     NOT NULL DEFAULT 0,
    reset_at    TIMESTAMPTZ NOT NULL
);
