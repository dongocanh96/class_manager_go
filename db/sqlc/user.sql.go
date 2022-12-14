// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    fullname,
    email,
    phone_number,
    is_teacher
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher
`

type CreateUserParams struct {
	Username       sql.NullString `json:"username"`
	HashedPassword string         `json:"hashed_password"`
	Fullname       sql.NullString `json:"fullname"`
	Email          sql.NullString `json:"email"`
	PhoneNumber    sql.NullString `json:"phone_number"`
	IsTeacher      bool           `json:"is_teacher"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.Fullname,
		arg.Email,
		arg.PhoneNumber,
		arg.IsTeacher,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PhoneNumber,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsTeacher,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getByUsername = `-- name: GetByUsername :one
SELECT id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetByUsername(ctx context.Context, username sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PhoneNumber,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsTeacher,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PhoneNumber,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsTeacher,
	)
	return i, err
}

const getUserForUpdate = `-- name: GetUserForUpdate :one
SELECT id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher FROM users
WHERE id = $1 LIMIT 1
FOR UPDATE
`

func (q *Queries) GetUserForUpdate(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserForUpdate, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PhoneNumber,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsTeacher,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher FROM users
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.HashedPassword,
			&i.Fullname,
			&i.Email,
			&i.PhoneNumber,
			&i.PasswordChangedAt,
			&i.CreatedAt,
			&i.IsTeacher,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserInfo = `-- name: UpdateUserInfo :one
UPDATE users
SET username = COALESCE($2, username),
    fullname = COALESCE($3, fullname),
    email = COALESCE($4, email),
    phone_number = COALESCE($5, phone_number)
WHERE id = $1
RETURNING id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher
`

type UpdateUserInfoParams struct {
	ID          int64          `json:"id"`
	Username    sql.NullString `json:"username"`
	Fullname    sql.NullString `json:"fullname"`
	Email       sql.NullString `json:"email"`
	PhoneNumber sql.NullString `json:"phone_number"`
}

func (q *Queries) UpdateUserInfo(ctx context.Context, arg UpdateUserInfoParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserInfo,
		arg.ID,
		arg.Username,
		arg.Fullname,
		arg.Email,
		arg.PhoneNumber,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PhoneNumber,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsTeacher,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2,
    password_changed_at = $3
WHERE id = $1
RETURNING id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher
`

type UpdateUserPasswordParams struct {
	ID                int64     `json:"id"`
	HashedPassword    string    `json:"hashed_password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.ID, arg.HashedPassword, arg.PasswordChangedAt)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.Fullname,
		&i.Email,
		&i.PhoneNumber,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsTeacher,
	)
	return i, err
}
