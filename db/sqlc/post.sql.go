// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: post.sql

package db

import (
	"context"
	"time"

	"github.com/lib/pq"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (title, body, user_id)
VALUES ($1, $2, $3) RETURNING id, title, body, user_id, created_at
`

type CreatePostParams struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int32  `json:"user_id"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost, arg.Title, arg.Body, arg.UserID)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Body,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE
FROM posts
WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPostById = `-- name: GetPostById :one
SELECT id, title, body, user_id, created_at
FROM posts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPostById(ctx context.Context, id int32) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPostById, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Body,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

const getPostsAndPostAuthors = `-- name: GetPostsAndPostAuthors :many
SELECT p.id, p.title, p.body, p.user_id, p.created_at,
       EXISTS(SELECT user_id, post_id FROM likes l WHERE l.post_id = p.id AND l.user_id = $1) AS is_liked,
       lc.likes_count,
       cc.comments_count,
       ARRAY(SELECT p_i.image_url FROM post_images p_i WHERE p_i.post_id = p.id)::text[] AS image_urls,
        u.id AS user_id,
       u.first_name,
       u.last_name,
       u.avatar_url,
       u.bio
FROM (SELECT id, title, body, user_id, created_at FROM posts ORDER BY created_at DESC OFFSET $3 LIMIT $2) p
         LEFT JOIN (SELECT post_id, COUNT(*) AS likes_count FROM likes GROUP BY post_id) lc ON lc.post_id = p.id
         LEFT JOIN (SELECT post_id, COUNT(*) AS comments_count FROM comments GROUP BY post_id) cc ON cc.post_id = p.id
         LEFT JOIN post_images p_i ON p_i.post_id = p.id
         JOIN users u ON p.user_id = u.id
`

type GetPostsAndPostAuthorsParams struct {
	UserID int32 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetPostsAndPostAuthorsRow struct {
	ID            int32     `json:"id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	UserID        int32     `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
	IsLiked       bool      `json:"is_liked"`
	LikesCount    int64     `json:"likes_count"`
	CommentsCount int64     `json:"comments_count"`
	ImageUrls     []string  `json:"image_urls"`
	UserID_2      int32     `json:"user_id_2"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	AvatarUrl     string    `json:"avatar_url"`
	Bio           string    `json:"bio"`
}

func (q *Queries) GetPostsAndPostAuthors(ctx context.Context, arg GetPostsAndPostAuthorsParams) ([]GetPostsAndPostAuthorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsAndPostAuthors, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPostsAndPostAuthorsRow{}
	for rows.Next() {
		var i GetPostsAndPostAuthorsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Body,
			&i.UserID,
			&i.CreatedAt,
			&i.IsLiked,
			&i.LikesCount,
			&i.CommentsCount,
			pq.Array(&i.ImageUrls),
			&i.UserID_2,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.Bio,
		); err != nil {
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

const listPosts = `-- name: ListPosts :many
SELECT id, title, body, user_id, created_at
FROM posts
ORDER BY created_at
`

func (q *Queries) ListPosts(ctx context.Context) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, listPosts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Body,
			&i.UserID,
			&i.CreatedAt,
		); err != nil {
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
