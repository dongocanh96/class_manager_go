-- name: CreateUser :one
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

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR UPDATE;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListTeachersOrStudents :many
SELECT * FROM users
WHERE is_teacher = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateUserInfo :one
UPDATE users
SET username = COALESCE($2, username),
    fullname = COALESCE($3, fullname),
    email = COALESCE($4, email),
    phone_number = COALESCE($5, phone_number)
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2,
    password_changed_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
