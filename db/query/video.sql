-- name: VideoSelectVideos :many
SELECT * FROM videos;

-- name: VideoSelectVideo :one
SELECT * FROM videos where slug=?1;

-- name: VideoSelectVideosByChannel :many
SELECT * FROM videos where uploader=?1;

-- name: VideoInsertVideo :one
INSERT INTO
	   videos(slug, title, description, uploader)
	   values(@slug, @title, @description, @uploader)
RETURNING *;

-- name: VideoDeleteVideo :one
DELETE FROM
	   videos
	   WHERE slug=?1
RETURNING *;

-- name: VideoUpdateVideo :one
UPDATE videos
	   SET title=@title, description=@description, processed=@processed
	   WHERE slug=@slug
RETURNING *;
