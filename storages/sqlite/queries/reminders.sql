-- name: GetReminders :many
SELECT
    *
FROM
    reminder;

-- name: AddReminder :exec
INSERT INTO
    reminder (message_id, user_id, channel_id, remind_at)
VALUES
    (?, ?, ?, ?);

-- name: DeleteReminder :exec
DELETE FROM
    reminder
WHERE
    id = ?;