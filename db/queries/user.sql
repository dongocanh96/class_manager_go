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

-- name: UpdateUser :one
UPDATE users
SET username = COALESCE($2, username),
    hashed_password = COALESCE($3, hashed_password),
    password_changed_at = COALESCE($4, password_changed_at),
    fullname = COALESCE($5, fullname),
    email = COALESCE($6, email),
    phone_number = COALESCE($7, phone_number)
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
