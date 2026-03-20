-- +goose Up
-- +goose StatementBegin
-- To generate UUIDv7 we need to install pg_uuidv7
CREATE EXTENSION IF NOT EXISTS "pg_uuidv7";

CREATE TABLE urls {
    id              UUID            PRIMARY KEY DEFAULT uuid_generate_v7(),
    code            TEXT            NOT NULL UNIQUE,
    url             TEXT            NOT NULL,
    account_number  VARCHAR(10)     ON DELETE SET NULL,
    expired_at      TIMESTAMPZ,
    click_count     BIGINT          NOT NULL DEFAULT 0,
    state           TEXT            NOT NULL DEFAULT 'active'
                    CHECK (state IN ('active', 'inactive', 'expired', 'disabled', 'archived', 'pending')),
    metadata        JSONB,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
};

-- Indexes
CREATE INDEX idx_urls_account_number
    ON urls(account_number);

CREATE INDEX idx_urls_expired_at
    on url(expired_at)
    WHERE deleted_at IS NOT NULL;   -- only index active (non-deleted) rows

-- Auto-update updated_at trigger
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER urls_set_updated_at
    BEFORE UPDATE ON urls
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS urls_set_updated_at ON urls;
DROP FUNCTION IF EXISTS set_updated_at;
DROP TABLE urls;
-- +goose StatementEnd
