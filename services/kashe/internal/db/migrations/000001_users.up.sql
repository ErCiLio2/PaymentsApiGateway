CREATE TABLE users (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email             TEXT        UNIQUE NOT NULL,
    password_hash     TEXT,
    google_id         TEXT        UNIQUE,
    name              TEXT        NOT NULL,
    email_verified_at TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
