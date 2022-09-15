-- +goose Up
CREATE TABLE reminder (
    id INTEGER PRIMARY KEY,
    message_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    channel_id TEXT NOT NULL,
    remind_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE reminder;