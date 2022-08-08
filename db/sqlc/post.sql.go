// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: post.sql

package db

import (
	"context"
	"time"
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
SELECT p.id AS post_id, p.title, p.body, p.created_at, u.id AS user_id, u.first_name, u.last_name
FROM posts p
         JOIN users u ON p.user_id = u.id
ORDER BY created_at LIMIT $1
OFFSET $2
`

type GetPostsAndPostAuthorsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetPostsAndPostAuthorsRow struct {
	PostID    int32     `json:"post_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UserID    int32     `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

func (q *Queries) GetPostsAndPostAuthors(ctx context.Context, arg GetPostsAndPostAuthorsParams) ([]GetPostsAndPostAuthorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPostsAndPostAuthors, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPostsAndPostAuthorsRow{}
	for rows.Next() {
		var i GetPostsAndPostAuthorsRow
		if err := rows.Scan(
			&i.PostID,
			&i.Title,
			&i.Body,
			&i.CreatedAt,
			&i.UserID,
			&i.FirstName,
			&i.LastName,
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
