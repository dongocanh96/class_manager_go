-- name: CreateMessage :one
INSERT INTO messages (
  from_user_id,
  to_user_id,
  content
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1 LIMIT 1;

-- name: ListMessages :many
SELECT * FROM messages
WHERE 
    from_user_id = $1 OR
    to_user_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: ListMessagesFromUser :many
SELECT * FROM messages
WHERE from_user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListMessagesToUser :many
SELECT * FROM messages
WHERE to_user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateMessage :one
UPDATE messages
SET content = $2
WHERE id = $1
RETURNING *;

-- name: UpdateMessageState :one
UPDATE messages
SET is_read = $2,
    read_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1;
