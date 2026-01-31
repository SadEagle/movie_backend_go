-- name: GetUserFavoriteList :many
SELECT movie_id
FROM favorite
WHERE user_id = $1;

-- name: GetMovieFavoriteList :many
SELECT user_id
FROM favorite
WHERE movie_id = $1;

-- name: GetFavorite :one
SELECT *
FROM favorite
WHERE movie_id = $1 
  AND user_id = $2;

-- name: CreateFavorite :one
INSERT INTO favorite(user_id, movie_id)
VALUES ($1, $2)
RETURNING user_id, movie_id;

-- name: DeleteFavorite :execrows
DELETE FROM favorite
WHERE user_id = $1 
  AND movie_id = $2;
