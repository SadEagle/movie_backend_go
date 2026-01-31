-- name: GetMovieCommentList :many
SELECT id, user_id, text, created_at
FROM comment
WHERE movie_id = $1;

-- name: GetUserCommentList :many
SELECT id, movie_id, text, created_at
FROM comment
WHERE user_id = $1;

-- name: GetComment :one
SELECT *
FROM comment
WHERE id = $1;

-- name: CreateComment :one
INSERT INTO comment (user_id, movie_id, text)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateComment :one
UPDATE comment SET
  text = $2
WHERE id = $1 
RETURNING *;


-- name: DeleteComment :execrows
DELETE FROM comment
WHERE id = $1
