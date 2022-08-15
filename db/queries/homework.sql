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
SET (
    title,
    file_name,
    saved_path,
    is_closed,
    updated_at,
    closed_at
) VALUES (
    $2, $3, $4, $5, $6, $7
)
WHERE id = $1
RETURNING *;

-- name: DeleteHomework :exec
DELETE FROM homeworks
WHERE id = $1;
