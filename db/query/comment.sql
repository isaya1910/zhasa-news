-- Example queries for sqlc

-- name: GetCommentById :one
SELECT *
FROM comments
WHERE id = $1 LIMIT 1;

-- name: GetCommentsAndAuthorsByPostId :many
SELECT comments(id) as comment_id, body, user_id, post_id, created_at, first_name, last_name
FROM comments
JOIN users
ON comments.user_id = users.id
WHERE post_id = $1
ORDER BY created_at;

-- name: CreateComment :one
INSERT INTO comments (body, user_id, post_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteComment :exec
DELETE
FROM comments
WHERE id = $1;
