-- name: GetMovieCommentList :many
SELECT id, user_id, text, created_at
FROM movie_comment
WHERE movie_id = $1;

-- name: GetUserCommentList :many
SELECT id, movie_id, text, created_at
FROM movie_comment
WHERE user_id = $1;


-- name: CreateMovieComment :one
INSERT INTO movie_comment (user_id, movie_id, text)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateMovieComment :one
UPDATE movie_comment SET
  text = $1
WHERE user_id = $2 
  AND movie_id = $3
RETURNING *;


-- name: DeleteMovieComment :execrows
DELETE FROM movie_comment
WHERE user_id = $1 
  AND movie_id = $2;
