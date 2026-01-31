-- name: GetMovieByID :one
SELECT id, title, COALESCE(amount_rates, 0), COALESCE(rating, 0), created_at
FROM (
  select * from movie where id = $1
  ) m
LEFT JOIN ( 
  select * from total_rating_mview where movie_id = $1
) mrv ON m.id = mrv.movie_id;

-- name: GetMovieByTitle :one
SELECT id, title, COALESCE(amount_rates, 0), COALESCE(rating, 0), created_at
FROM (
  select * from movie where title = $1
  ) m
LEFT JOIN total_rating_mview ON m.id = mrv.movie_id;

-- name: GetMovieList :many
SELECT *
FROM movie;

-- name: CreateMovie :one
INSERT INTO movie(title)
VALUES ($1)
RETURNING *;

-- name: UpdateMovie :one
UPDATE movie SET
  title = COALESCE(sqlc.narg(title), title)
WHERE id = $1
RETURNING *;

-- name: DeleteMovie :execrows
DELETE FROM movie
WHERE id = $1;
