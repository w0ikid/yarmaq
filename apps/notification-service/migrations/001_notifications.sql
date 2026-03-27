-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notifications (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    VARCHAR(255) NOT NULL,
    type       VARCHAR(100) NOT NULL,
    channel    VARCHAR(20) NOT NULL,
    status     VARCHAR(20) NOT NULL,
    subject    TEXT NOT NULL,
    body       TEXT NOT NULL,
    metadata   JSONB,
    error      TEXT,
    sent_at    TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
-- +goose StatementEnd
