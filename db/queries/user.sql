-- name: CreateTeacherUser :one
INSERT INTO users (
    username,
    hashed_password,
    fullname,
    email,
    phone_number,
    is_teacher
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: CreateStudentUser :one
INSERT INTO users (
    username,
    hashed_password,
    fullname,
    email,
    phone_number
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET username = $2,
    hashed_password = $3,
    fullname = $4,
    email = $5,
    phone_number = $6,
    password_changed_at = $7
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
