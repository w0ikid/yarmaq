-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions
    ADD COLUMN target_amount BIGINT,
    ADD COLUMN target_currency VARCHAR(3),
    ADD COLUMN exchange_rate NUMERIC;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
    DROP COLUMN IF EXISTS target_amount,
    DROP COLUMN IF EXISTS target_currency,
    DROP COLUMN IF EXISTS exchange_rate;
-- +goose StatementEnd
