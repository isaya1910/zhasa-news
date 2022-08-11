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
SELECT p.*,
       EXISTS(SELECT * FROM likes l WHERE l.post_id = p.id AND l.user_id = $1) AS is_liked,
       (SELECT COUNT(*) AS likes_count FROM likes l WHERE l.post_id = p.id),
       u.first_name,
       u.last_name
FROM posts p
         LEFT JOIN likes l ON l.post_id = p.id
         JOIN users u ON p.user_id = u.id
ORDER BY created_at LIMIT $2
OFFSET $3;