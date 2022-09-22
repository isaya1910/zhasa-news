-- Example queries for sqlc

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY name;

-- name: CreateOrUpdateUser :one
INSERT INTO users (id, first_name, last_name, bio, avatar_url)
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id)
DO
UPDATE
    SET first_name = excluded.first_name,
    last_name = excluded.last_name,
    bio = excluded.bio,
    avatar_url = excluded.avatar_url
    RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;
