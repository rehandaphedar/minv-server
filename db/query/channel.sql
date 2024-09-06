-- name: ChannelSelectChannel :one
SELECT
    channelname, created
FROM
    channels
WHERE
    channelname = ?1;

-- name: ChannelSelectChannels :many
SELECT
    channelname, created
FROM
    channels;

-- name: ChannelDeleteChannel :one
DELETE FROM
	   channels
WHERE
	channelname = ?1
RETURNING channelname, created;
