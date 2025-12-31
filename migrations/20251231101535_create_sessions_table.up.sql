CREATE TABLE sessions (
    token varchar NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX idx_sessions_token
ON sessions (token);

CREATE INDEX idx_sessions_expires_at
ON sessions (expires_at);

CREATE INDEX idx_sessions_user_id
ON sessions (user_id);

