// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: image.sql

package db

import (
	"context"
)

const createPostImage = `-- name: CreatePostImage :one
INSERT INTO post_images(image_url, post_id)
VALUES (
    $1, $2
)
RETURNING id, image_url, post_id
`

type CreatePostImageParams struct {
	ImageUrl string `json:"image_url"`
	PostID   int32  `json:"post_id"`
}

func (q *Queries) CreatePostImage(ctx context.Context, arg CreatePostImageParams) (PostImage, error) {
	row := q.db.QueryRowContext(ctx, createPostImage, arg.ImageUrl, arg.PostID)
	var i PostImage
	err := row.Scan(&i.ID, &i.ImageUrl, &i.PostID)
	return i, err
}

const deleteImage = `-- name: DeleteImage :exec
DELETE FROM post_images
WHERE post_id = $1
`

func (q *Queries) DeleteImage(ctx context.Context, postID int32) error {
	_, err := q.db.ExecContext(ctx, deleteImage, postID)
	return err
}

const getPostImages = `-- name: GetPostImages :many
SELECT id, image_url, post_id
from post_images
WHERE post_id = $1
`

func (q *Queries) GetPostImages(ctx context.Context, postID int32) ([]PostImage, error) {
	rows, err := q.db.QueryContext(ctx, getPostImages, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PostImage{}
	for rows.Next() {
		var i PostImage
		if err := rows.Scan(&i.ID, &i.ImageUrl, &i.PostID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
