-- name: GetUserByID :one
SELECT id, name, login, encoded_password, is_admin, created_at
FROM user_data
WHERE id = $1;

-- name: GetUserByLogin :one
SELECT id, name, login, encoded_password, is_admin, created_at
FROM user_data
WHERE login = $1;

-- name: GetUserList :many
SELECT id, name, login, encoded_password, is_admin, created_at
FROM user_data;

-- name: CreateUser :one
INSERT INTO user_data(name, login, encoded_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE user_data SET
  name = COALESCE(sqlc.narg(name), name),
  login = COALESCE(sqlc.narg(login), login),
  encoded_password = COALESCE(sqlc.narg(encode_password), encode_password)
WHERE id = $1
RETURNING *;

-- name: DeleteUser :execrows
DELETE FROM user_data
WHERE id = $1;
