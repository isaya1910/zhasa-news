-- Example queries for sqlc

-- name: GetPostById :one
SELECT *
FROM posts
WHERE id = $1 LIMIT 1;

-- name: ListPosts :many
SELECT *
FROM posts
ORDER BY created_at;

-- name: CreatePost :one
INSERT INTO posts (title, body, user_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: DeletePost :exec
DELETE
FROM posts
WHERE id = $1;

-- name: GetPostsAndPostAuthors :many
SELECT p.id AS post_id, p.title, p.body, p.created_at, u.id AS user_id, u.first_name, u.last_name
FROM posts p
         JOIN users u ON p.user_id = u.id
ORDER BY created_at LIMIT $1
OFFSET $2;