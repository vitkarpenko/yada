-- name: GetMuseRating :one
SELECT
    rating
FROM
    muse
WHERE
    hash = ?
LIMIT
    1;

-- name: CreateMuse :exec
INSERT INTO
    muse (hash, rating)
VALUES
    (?, ?);