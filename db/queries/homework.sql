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

-- name: ListHomeworksByTeacher :many
SELECT * FROM homeworks
WHERE teacher_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListHomeworksBySubject :many
SELECT * FROM homeworks
WHERE subject = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListHomeworks :many
SELECT * FROM homeworks
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateHomework :one
UPDATE homeworks
SET title = COALESCE($2, title),
    file_name = COALESCE($3, file_name),
    saved_path = COALESCE($4, saved_path),
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
