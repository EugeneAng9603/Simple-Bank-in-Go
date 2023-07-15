-- name: CreateUser :one
INSERT INTO users (
  username_of_user,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username_of_user = $1 LIMIT 1;