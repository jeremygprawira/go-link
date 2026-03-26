-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id              UUID            PRIMARY KEY,
    code            TEXT            NOT NULL UNIQUE,
    name            TEXT            NOT NULL,
    url             TEXT            NOT NULL,
    account_number  VARCHAR(10)     NOT NULL,
    click_count     BIGINT          NOT NULL DEFAULT 0,
    state           TEXT            NOT NULL DEFAULT 'active'
                    CHECK (state IN ('active', 'inactive', 'expired', 'disabled', 'archived', 'pending')),
    metadata        JSONB,
    expired_at      TIMESTAMPTZ     DEFAULT NULL,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ     DEFAULT NULL
);

-- Indexes
CREATE INDEX idx_urls_account_number
    ON urls(account_number);

CREATE INDEX idx_urls_url
    ON urls(url);

CREATE INDEX idx_urls_expired_at
    on urls(expired_at)
    WHERE deleted_at IS NULL;   -- only index active (non-deleted) rows

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
DROP INDEX IF EXISTS idx_urls_account_number;
DROP INDEX IF EXISTS idx_urls_url;
DROP INDEX IF EXISTS idx_urls_expired_at;
DROP TRIGGER IF EXISTS urls_set_updated_at ON urls;
DROP FUNCTION IF EXISTS set_updated_at;
DROP TABLE urls;
-- +goose StatementEnd
