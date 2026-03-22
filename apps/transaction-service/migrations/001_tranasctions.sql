-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type            VARCHAR(20) NOT NULL DEFAULT 'TRANSFER',
    from_account_id UUID NOT NULL,
    to_account_id   UUID NOT NULL,
    amount          BIGINT NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'KZT',
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    -- PENDING, HOLDING, DEPOSITING, COMPLETED, FAILED
    idempotency_key VARCHAR(255) UNIQUE,  -- защита от дублей
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
