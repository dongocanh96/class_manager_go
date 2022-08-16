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

-- name: UpdateUsername :one
UPDATE users
SET username = $2
WHERE id = $1
RETURNING *;

-- name: UpdateHashedPassword :one
UPDATE users
SET hashed_password = $2,
    password_changed_at = $3
WHERE id = $1
RETURNING *;

-- name: UpdateFullname :one
UPDATE users
SET fullname = $2
WHERE id = $1
RETURNING *;

-- name: UpdateEmail :one
UPDATE users
SET email = $2
WHERE id = $1
RETURNING *;

-- name: UpdatePhoneNumber :one
UPDATE users
SET phone_number = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
