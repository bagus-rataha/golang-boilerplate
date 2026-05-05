CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         BIGSERIAL PRIMARY KEY,
    token      TEXT NOT NULL,
    user_id    BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expired_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens (token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expired_at ON refresh_tokens (expired_at);
