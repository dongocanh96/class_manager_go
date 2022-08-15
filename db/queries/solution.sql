-- name: CreateSolution :one
INSERT INTO solutions (
    problem_id,
    user_id,
    file_name,
    saved_path
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetSolution :one
SELECT * FROM solutions
WHERE problem_id = $1 AND user_id = $2
LIMIT 1;

-- name: ListSolutionsByUser :many
SELECT * FROM solutions
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;


-- name: ListSolutionsByProblem :many
SELECT * FROM solutions
WHERE problem_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: ListSolutions :many
SELECT * FROM solutions
WHERE 
    user_id = $1 OR
    problem_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateSolution :one
UPDATE solutions
SET (
    file_name,
    saved_path,
    updated_at,
) VALUES (
    $2, $3, $4
)
WHERE user_id = $1
RETURNING *;

-- name: DeleteSolution :exec
DELETE FROM solutions
WHERE id = $1;
