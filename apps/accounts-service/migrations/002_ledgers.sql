-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ledgers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id      UUID NOT NULL REFERENCES accounts(id),
    amount          BIGINT NOT NULL,          -- 
    operation_type  VARCHAR(50) NOT NULL,    -- 'DEPOSIT', 'WITHDRAW', 'HOLD', 'REFUND'
    reference_id    UUID,                     -- ID транзакции из Saga
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_ledger_account_created ON ledgers(account_id, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ledgers;
-- +goose StatementEnd