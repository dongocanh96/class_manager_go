-- name: CreateHomework :one
INSERT INTO homeworks (
    teacher_id,
    subject,
    title,
    file_name,
    saved_path
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetHomework :one
SELECT * FROM homeworks
WHERE id = $1 LIMIT 1;

-- name: ListHomeworkByTeacherId :many
SELECT * FROM homeworks
WHERE teacher_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListHomework :many
SELECT * FROM homeworks
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateHomework :one
UPDATE homeworks
SET title = $2,
    file_name = $3,
    saved_path = $4,
    updated_at = $5
WHERE id = $1
RETURNING *;

-- name: CloseHomework :one
UPDATE homeworks
SET is_closed = $2,
    closed_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteHomework :exec
DELETE FROM homeworks
WHERE id = $1;
