-- +goose Up
CREATE TABLE muse (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hash TEXT NOT NULL,
    rating INTEGER NOT NULL
);

-- +goose Down
DROP TABLE muse;