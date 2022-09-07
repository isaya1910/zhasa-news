-- name: GetPostImages :many
SELECT *
from post_images
WHERE post_id = $1;

-- name: CreatePostImage :one
INSERT INTO post_images(image_url, post_id)
VALUES (
    $1, $2
)
RETURNING *;

-- name: DeleteImage :exec
DELETE FROM post_images
WHERE post_id = $1;