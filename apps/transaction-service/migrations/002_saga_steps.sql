-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS saga_steps (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id  UUID NOT NULL REFERENCES transactions(id),
    step_name       VARCHAR(50) NOT NULL,
    -- 'HOLD', 'DEPOSIT', 'WITHDRAW', 'REFUND'
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    -- PENDING, COMPLETED, FAILED
    error           TEXT,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX idx_saga_steps_transaction ON saga_steps(transaction_id);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS saga_steps;
-- +goose StatementEnd