// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: message.sql

package db

import (
	"context"
	"time"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (
  from_user_id,
  to_user_id,
  title,
  content
) VALUES (
  $1, $2, $3, $4
) RETURNING id, from_user_id, to_user_id, title, content, is_read, created_at, read_at
`

type CreateMessageParams struct {
	FromUserID int64  `json:"from_user_id"`
	ToUserID   int64  `json:"to_user_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, createMessage,
		arg.FromUserID,
		arg.ToUserID,
		arg.Title,
		arg.Content,
	)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.Title,
		&i.Content,
		&i.IsRead,
		&i.CreatedAt,
		&i.ReadAt,
	)
	return i, err
}

const deleteMessage = `-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1
`

func (q *Queries) DeleteMessage(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteMessage, id)
	return err
}

const getMessage = `-- name: GetMessage :one
SELECT id, from_user_id, to_user_id, title, content, is_read, created_at, read_at FROM messages
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetMessage(ctx context.Context, id int64) (Message, error) {
	row := q.db.QueryRowContext(ctx, getMessage, id)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.Title,
		&i.Content,
		&i.IsRead,
		&i.CreatedAt,
		&i.ReadAt,
	)
	return i, err
}

const listMessages = `-- name: ListMessages :many
SELECT id, from_user_id, to_user_id, title, content, is_read, created_at, read_at FROM messages
WHERE 
    from_user_id = $1 OR
    to_user_id = $2
ORDER BY id
LIMIT $3
OFFSET $4
`

type ListMessagesParams struct {
	FromUserID int64 `json:"from_user_id"`
	ToUserID   int64 `json:"to_user_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListMessages(ctx context.Context, arg ListMessagesParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, listMessages,
		arg.FromUserID,
		arg.ToUserID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Message{}
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.FromUserID,
			&i.ToUserID,
			&i.Title,
			&i.Content,
			&i.IsRead,
			&i.CreatedAt,
			&i.ReadAt,
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

const listMessagesFromUser = `-- name: ListMessagesFromUser :many
SELECT id, from_user_id, to_user_id, title, content, is_read, created_at, read_at FROM messages
WHERE from_user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListMessagesFromUserParams struct {
	FromUserID int64 `json:"from_user_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListMessagesFromUser(ctx context.Context, arg ListMessagesFromUserParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, listMessagesFromUser, arg.FromUserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Message{}
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.FromUserID,
			&i.ToUserID,
			&i.Title,
			&i.Content,
			&i.IsRead,
			&i.CreatedAt,
			&i.ReadAt,
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

const listMessagesToUser = `-- name: ListMessagesToUser :many
SELECT id, from_user_id, to_user_id, title, content, is_read, created_at, read_at FROM messages
WHERE to_user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListMessagesToUserParams struct {
	ToUserID int64 `json:"to_user_id"`
	Limit    int32 `json:"limit"`
	Offset   int32 `json:"offset"`
}

func (q *Queries) ListMessagesToUser(ctx context.Context, arg ListMessagesToUserParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, listMessagesToUser, arg.ToUserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Message{}
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.FromUserID,
			&i.ToUserID,
			&i.Title,
			&i.Content,
			&i.IsRead,
			&i.CreatedAt,
			&i.ReadAt,
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

const updateMessage = `-- name: UpdateMessage :one
UPDATE messages
SET title = $2,
    content = $3,
    is_read = $4,
    read_at = $5
WHERE id = $1
RETURNING id, from_user_id, to_user_id, title, content, is_read, created_at, read_at
`

type UpdateMessageParams struct {
	ID      int64     `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	IsRead  bool      `json:"is_read"`
	ReadAt  time.Time `json:"read_at"`
}

func (q *Queries) UpdateMessage(ctx context.Context, arg UpdateMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, updateMessage,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.IsRead,
		arg.ReadAt,
	)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.FromUserID,
		&i.ToUserID,
		&i.Title,
		&i.Content,
		&i.IsRead,
		&i.CreatedAt,
		&i.ReadAt,
	)
	return i, err
}
