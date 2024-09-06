-- name: AuthSelectChannel :one
SELECT
    *
FROM
    channels
WHERE
    channelname = ?1;

-- name: AuthInsertChannel :one
INSERT INTO channels (channelname, password)
    VALUES (?1, ?2)
RETURNING
    *;

