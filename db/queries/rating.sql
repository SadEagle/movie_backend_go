-- name: GetUserRatingList :many
SELECT movie_id, rating
FROM rating
WHERE user_id = $1;

-- name: GetMovieRatingList :many
SELECT movie_id, rating
FROM rating
WHERE user_id = $1;

-- name: GetRating :one
SELECT *
FROM rating
WHERE user_id = $1 
  AND movie_id = $2;

-- name: CreateRating :one
INSERT INTO rating (user_id, movie_id, rating)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateRating :one
UPDATE rating SET
  rating = COALESCE(sqlc.narg(rating), rating)
WHERE user_id = $1 
  AND movie_id = $2
RETURNING *;


-- name: DeleteRating :execrows
DELETE FROM rating
WHERE user_id = $1 
  AND movie_id = $2;
