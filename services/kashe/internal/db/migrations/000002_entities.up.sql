CREATE TABLE entities (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name           TEXT        NOT NULL,
    description    TEXT,
    payment_method TEXT        NOT NULL DEFAULT 'proxypay',
    otp_method     TEXT        NOT NULL DEFAULT 'email',
    totp_secret    TEXT,
    token_hash     TEXT        UNIQUE NOT NULL,
    token_prefix   TEXT        NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON entities(user_id);
CREATE INDEX ON entities(token_hash);
