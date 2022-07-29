// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: comment.sql

package db

import (
	"context"
	"time"
)

const createComment = `-- name: CreateComment :one
INSERT INTO comments (body, user_id, post_id)
VALUES ($1, $2, $3) RETURNING id, body, user_id, post_id, created_at
`

type CreateCommentParams struct {
	Body   string `json:"body"`
	UserID int32  `json:"user_id"`
	PostID int32  `json:"post_id"`
}

func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRowContext(ctx, createComment, arg.Body, arg.UserID, arg.PostID)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Body,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteComment = `-- name: DeleteComment :exec
DELETE
FROM comments
WHERE id = $1
`

func (q *Queries) DeleteComment(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteComment, id)
	return err
}

const getCommentById = `-- name: GetCommentById :one

SELECT id, body, user_id, post_id, created_at
FROM comments
WHERE id = $1 LIMIT 1
`

// Example queries for sqlc
func (q *Queries) GetCommentById(ctx context.Context, id int32) (Comment, error) {
	row := q.db.QueryRowContext(ctx, getCommentById, id)
	var i Comment
	err := row.Scan(
		&i.ID,
		&i.Body,
		&i.UserID,
		&i.PostID,
		&i.CreatedAt,
	)
	return i, err
}

const getCommentsAndAuthorsByPostId = `-- name: GetCommentsAndAuthorsByPostId :many
SELECT c.id as comment_id, c.body, c.user_id, c.post_id, c.created_at, u.first_name, u.last_name
FROM comments c
JOIN users u
ON c.user_id = u.id
WHERE c.post_id = $1
ORDER BY created_at
`

type GetCommentsAndAuthorsByPostIdRow struct {
	CommentID int32     `json:"comment_id"`
	Body      string    `json:"body"`
	UserID    int32     `json:"user_id"`
	PostID    int32     `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

func (q *Queries) GetCommentsAndAuthorsByPostId(ctx context.Context, postID int32) ([]GetCommentsAndAuthorsByPostIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getCommentsAndAuthorsByPostId, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCommentsAndAuthorsByPostIdRow{}
	for rows.Next() {
		var i GetCommentsAndAuthorsByPostIdRow
		if err := rows.Scan(
			&i.CommentID,
			&i.Body,
			&i.UserID,
			&i.PostID,
			&i.CreatedAt,
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
