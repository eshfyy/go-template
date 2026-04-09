-- +goose Up
CREATE TABLE notifications (
    id         UUID PRIMARY KEY,
    user_id    UUID        NOT NULL REFERENCES users(id),
    title      TEXT        NOT NULL,
    text       TEXT        NOT NULL,
    channel    TEXT        NOT NULL,
    status     TEXT        NOT NULL DEFAULT 'pending',
    sent_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_status_created ON notifications(status, created_at);

-- +goose Down
DROP TABLE notifications;
