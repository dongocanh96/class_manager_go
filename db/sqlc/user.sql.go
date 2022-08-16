// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: user.sql

package db

import (
	"context"
	"time"
)

const createStudentUser = `-- name: CreateStudentUser :one
INSERT INTO users (
    username,
    hashed_password,
    fullname,
    email,
    phone_number
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher
`

type CreateStudentUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
}

func (q *Queries) CreateStudentUser(ctx context.Context, arg CreateStudentUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createStudentUser,
		arg.Username,
		arg.HashedPassword,
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

const createTeacherUser = `-- name: CreateTeacherUser :one
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

type CreateTeacherUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	Fullname       string `json:"fullname"`
	Email          string `json:"email"`
	PhoneNumber    string `json:"phone_number"`
	IsTeacher      bool   `json:"is_teacher"`
}

func (q *Queries) CreateTeacherUser(ctx context.Context, arg CreateTeacherUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createTeacherUser,
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

const updateEmail = `-- name: UpdateEmail :one
UPDATE users
SET email = $2
WHERE id = $1
RETURNING id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher
`

type UpdateEmailParams struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func (q *Queries) UpdateEmail(ctx context.Context, arg UpdateEmailParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateEmail, arg.ID, arg.Email)
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

const updateFullname = `-- name: UpdateFullname :one
UPDATE users
SET fullname = $2
WHERE id = $1
RETURNING id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher
`

type UpdateFullnameParams struct {
	ID       int64  `json:"id"`
	Fullname string `json:"fullname"`
}

func (q *Queries) UpdateFullname(ctx context.Context, arg UpdateFullnameParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateFullname, arg.ID, arg.Fullname)
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

const updateHashedPassword = `-- name: UpdateHashedPassword :one
UPDATE users
SET hashed_password = $2
WHERE id = $1
RETURNING id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher
`

type UpdateHashedPasswordParams struct {
	ID             int64  `json:"id"`
	HashedPassword string `json:"hashed_password"`
}

func (q *Queries) UpdateHashedPassword(ctx context.Context, arg UpdateHashedPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateHashedPassword, arg.ID, arg.HashedPassword)
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

const updatePasswordChangedTime = `-- name: UpdatePasswordChangedTime :one
UPDATE users
SET password_changed_at = $2
WHERE id = $1
RETURNING id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher
`

type UpdatePasswordChangedTimeParams struct {
	ID                int64     `json:"id"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func (q *Queries) UpdatePasswordChangedTime(ctx context.Context, arg UpdatePasswordChangedTimeParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updatePasswordChangedTime, arg.ID, arg.PasswordChangedAt)
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

const updatePhoneNumber = `-- name: UpdatePhoneNumber :one
UPDATE users
SET phone_number = $2
WHERE id = $1
RETURNING id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher
`

type UpdatePhoneNumberParams struct {
	ID          int64  `json:"id"`
	PhoneNumber string `json:"phone_number"`
}

func (q *Queries) UpdatePhoneNumber(ctx context.Context, arg UpdatePhoneNumberParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updatePhoneNumber, arg.ID, arg.PhoneNumber)
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

const updateUserName = `-- name: UpdateUserName :one
UPDATE users
SET username = $2
WHERE id = $1
RETURNING id, username, hashed_password, fullname, email, phone_number, password_changed_at, created_at, is_teacher
`

type UpdateUserNameParams struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (q *Queries) UpdateUserName(ctx context.Context, arg UpdateUserNameParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserName, arg.ID, arg.Username)
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
