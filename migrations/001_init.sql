-- +goose Up
CREATE TABLE muse (
    hash TEXT PRIMARY KEY,
    rating INTEGER NOT NULL
);

-- +goose Down
DROP TABLE muse;