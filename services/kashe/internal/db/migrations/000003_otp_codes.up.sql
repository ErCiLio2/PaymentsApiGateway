CREATE TABLE otp_codes (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT        NOT NULL,
    code       TEXT        NOT NULL,
    purpose    TEXT        NOT NULL,
    charge_id  TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON otp_codes(email, purpose) WHERE used_at IS NULL;
