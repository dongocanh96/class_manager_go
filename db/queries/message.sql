-- name: CreateMessage :one
INSERT INTO messages (
  from_user_id,
  to_user_id,
  title,
  content
) VALUES (
  $1, $2, $3, $4
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
SET (
    title,
    content,
    is_read,
    read_at
) VALUES (
    $2, $3, $4, $5
)
WHERE id = $1
RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM transfers
WHERE id = $1;
