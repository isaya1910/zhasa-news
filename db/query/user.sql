-- Example queries for sqlc


SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: Listusers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO users (
  first_name, second_name 
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM users
WHERE id = $1;
