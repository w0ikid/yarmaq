-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type    VARCHAR(50) NOT NULL,
    payload       JSONB NOT NULL,
    aggregate_id  UUID NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    sent_at       TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_outbox_unsent ON outbox(created_at) WHERE sent_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox;
-- +goose StatementEnd