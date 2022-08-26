// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: solution.sql

package db

import (
	"context"
	"time"
)

const createSolution = `-- name: CreateSolution :one
INSERT INTO solutions (
    problem_id,
    user_id,
    file_name,
    saved_path
) VALUES (
    $1, $2, $3, $4
) RETURNING id, problem_id, user_id, file_name, saved_path, submited_at, updated_at
`

type CreateSolutionParams struct {
	ProblemID int64  `json:"problem_id"`
	UserID    int64  `json:"user_id"`
	FileName  string `json:"file_name"`
	SavedPath string `json:"saved_path"`
}

func (q *Queries) CreateSolution(ctx context.Context, arg CreateSolutionParams) (Solution, error) {
	row := q.db.QueryRowContext(ctx, createSolution,
		arg.ProblemID,
		arg.UserID,
		arg.FileName,
		arg.SavedPath,
	)
	var i Solution
	err := row.Scan(
		&i.ID,
		&i.ProblemID,
		&i.UserID,
		&i.FileName,
		&i.SavedPath,
		&i.SubmitedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteSolution = `-- name: DeleteSolution :exec
DELETE FROM solutions
WHERE id = $1
`

func (q *Queries) DeleteSolution(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteSolution, id)
	return err
}

const getSolutionByID = `-- name: GetSolutionByID :one
SELECT id, problem_id, user_id, file_name, saved_path, submited_at, updated_at FROM solutions
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetSolutionByID(ctx context.Context, id int64) (Solution, error) {
	row := q.db.QueryRowContext(ctx, getSolutionByID, id)
	var i Solution
	err := row.Scan(
		&i.ID,
		&i.ProblemID,
		&i.UserID,
		&i.FileName,
		&i.SavedPath,
		&i.SubmitedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSolutionByProblemAndUser = `-- name: GetSolutionByProblemAndUser :one
SELECT id, problem_id, user_id, file_name, saved_path, submited_at, updated_at FROM solutions
WHERE problem_id = $1 AND user_id = $2
LIMIT 1
`

type GetSolutionByProblemAndUserParams struct {
	ProblemID int64 `json:"problem_id"`
	UserID    int64 `json:"user_id"`
}

func (q *Queries) GetSolutionByProblemAndUser(ctx context.Context, arg GetSolutionByProblemAndUserParams) (Solution, error) {
	row := q.db.QueryRowContext(ctx, getSolutionByProblemAndUser, arg.ProblemID, arg.UserID)
	var i Solution
	err := row.Scan(
		&i.ID,
		&i.ProblemID,
		&i.UserID,
		&i.FileName,
		&i.SavedPath,
		&i.SubmitedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listSolutionsByProblem = `-- name: ListSolutionsByProblem :many
SELECT id, problem_id, user_id, file_name, saved_path, submited_at, updated_at FROM solutions
WHERE problem_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListSolutionsByProblemParams struct {
	ProblemID int64 `json:"problem_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListSolutionsByProblem(ctx context.Context, arg ListSolutionsByProblemParams) ([]Solution, error) {
	rows, err := q.db.QueryContext(ctx, listSolutionsByProblem, arg.ProblemID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Solution{}
	for rows.Next() {
		var i Solution
		if err := rows.Scan(
			&i.ID,
			&i.ProblemID,
			&i.UserID,
			&i.FileName,
			&i.SavedPath,
			&i.SubmitedAt,
			&i.UpdatedAt,
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

const listSolutionsByUser = `-- name: ListSolutionsByUser :many
SELECT id, problem_id, user_id, file_name, saved_path, submited_at, updated_at FROM solutions
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListSolutionsByUserParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListSolutionsByUser(ctx context.Context, arg ListSolutionsByUserParams) ([]Solution, error) {
	rows, err := q.db.QueryContext(ctx, listSolutionsByUser, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Solution{}
	for rows.Next() {
		var i Solution
		if err := rows.Scan(
			&i.ID,
			&i.ProblemID,
			&i.UserID,
			&i.FileName,
			&i.SavedPath,
			&i.SubmitedAt,
			&i.UpdatedAt,
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

const updateSolution = `-- name: UpdateSolution :one
UPDATE solutions
SET file_name = $2,
    saved_path = $3,
    updated_at = $4
WHERE id = $1
RETURNING id, problem_id, user_id, file_name, saved_path, submited_at, updated_at
`

type UpdateSolutionParams struct {
	ID        int64     `json:"id"`
	FileName  string    `json:"file_name"`
	SavedPath string    `json:"saved_path"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdateSolution(ctx context.Context, arg UpdateSolutionParams) (Solution, error) {
	row := q.db.QueryRowContext(ctx, updateSolution,
		arg.ID,
		arg.FileName,
		arg.SavedPath,
		arg.UpdatedAt,
	)
	var i Solution
	err := row.Scan(
		&i.ID,
		&i.ProblemID,
		&i.UserID,
		&i.FileName,
		&i.SavedPath,
		&i.SubmitedAt,
		&i.UpdatedAt,
	)
	return i, err
}
