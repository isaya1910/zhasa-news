-- Example queries for sqlc

-- name: GetCommentById :one
SELECT *
FROM comments
WHERE id = $1 LIMIT 1;

-- name: GetCommentsAndAuthorsByPostId :many
SELECT c.id as comment_id, c.body, c.user_id, c.post_id, c.created_at, u.id AS user_id, u.first_name, u.last_name, u.avatar_url, u.bio
FROM comments c
JOIN users u
ON c.user_id = u.id
WHERE c.post_id = $1
ORDER BY created_at;

-- name: CreateComment :one
INSERT INTO comments (body, user_id, post_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteComment :exec
DELETE
FROM comments
WHERE id = $1;
