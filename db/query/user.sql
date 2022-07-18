-- Example queries for sqlc

SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  id, first_name, second_name, bio 
) VALUES (
  $1, $2, $3, $4 
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
